package logger

import (
	"lambda-https-doctor/config"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var once sync.Once
var logger *zerolog.Logger

func Log() *zerolog.Logger {
	once.Do(func() {
		l := zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC1123,
		}).With().Timestamp().Logger().Level(config.Get().ZerologLevel())
		logger = &l
	})

	return logger
}
