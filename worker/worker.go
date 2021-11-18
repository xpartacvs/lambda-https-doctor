package worker

import (
	"errors"
	"lambda-https-doctor/alert"
	"lambda-https-doctor/client"
	"lambda-https-doctor/config"
	"lambda-https-doctor/logger"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
)

func Start() {
	lambda.Start(examine)
}

func examine() error {
	if len(config.Get().Hosts()) <= 0 {
		return errors.New("no host to check on")
	}

	notif := alert.New(config.Get().DishookBotMessage())
	notif.SetBotName(config.Get().DishookBotName()).SetBotAvatar(config.Get().DishookBotAvatar())

	var wg sync.WaitGroup
	for _, h := range config.Get().Hosts() {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			var title, content string = strings.ToUpper(host), ""
			httpsClient := client.New(host, logger.Log())
			expiry, err := httpsClient.GetExpiry()
			switch err {
			case client.ErrTimeout:
				content = "Connection timeout."
				if config.Get().ZerologLevel() != zerolog.Disabled {
					notif.AddField(title, content, false)
				}
			case client.ErrConnection:
				content = "Connection error."
				if config.Get().ZerologLevel() != zerolog.Disabled {
					notif.AddField(title, content, false)
				}
			case client.ErrCertInvalid:
				content = "Wrong SSL certificate."
				notif.AddField(title, content, false)
			case client.ErrCertExpired:
				days := int(time.Until(*expiry).Hours()/24) * -1
				content = "SSL expired for " + strconv.Itoa(days) + " days."
				notif.AddField(title, content, false)
			default:
				graceTime := expiry.AddDate(0, 0, config.Get().Graceperiod())
				if time.Now().After(graceTime) {
					days := int(time.Until(*expiry).Hours() / 24)
					content = "SSL expired in " + strconv.Itoa(days) + " days."
					notif.AddField(title, content, false)
				}
			}
		}(h)
	}
	wg.Wait()

	if err := notif.Send(config.Get().DishookURL(), true); err != nil {
		if err != alert.ErrNoContent {
			return err
		}
		logger.Log().Info().Msg(err.Error())
	}

	return nil
}
