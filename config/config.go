package config

import (
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type config struct {
	hosts          []string
	logLevel       zerolog.Level
	dishookBotMsg  string
	dishookBotName string
	dishookBotAva  string
	dishookUrl     string
	location       *time.Location
	grace          int
}

type Config interface {
	Hosts() []string
	ZerologLevel() zerolog.Level
	DishookBotMessage() string
	DishookBotName() string
	DishookBotAvatar() string
	DishookURL() string
	Location() *time.Location
	Graceperiod() int
}

var (
	cfg  *config
	once sync.Once
)

func (c *config) Hosts() []string {
	return c.hosts
}

func (c *config) ZerologLevel() zerolog.Level {
	return c.logLevel
}

func (c *config) DishookBotMessage() string {
	return c.dishookBotMsg
}

func (c *config) DishookBotName() string {
	return c.dishookBotName
}

func (c *config) DishookBotAvatar() string {
	return c.dishookBotAva
}

func (c *config) DishookURL() string {
	return c.dishookUrl
}

func (c *config) Location() *time.Location {
	return c.location
}

func (c *config) Graceperiod() int {
	return c.grace
}

func load() *config {
	fang := viper.New()

	fang.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	fang.AutomaticEnv()

	fang.SetConfigName("https-doctor")
	fang.SetConfigType("yml")
	fang.AddConfigPath(".")

	value, available := os.LookupEnv("CONFIG_LOCATION")
	if available {
		fang.AddConfigPath(value)
	}

	_ = fang.ReadInConfig()

	graceperiod := fang.GetInt("graceperiod")
	switch {
	case graceperiod == 0:
		graceperiod = -14
	case graceperiod > 0:
		graceperiod *= -1
	}

	return &config{
		hosts:          splitCSV(fang.GetString("hosts")),
		logLevel:       setLogLevel(fang.GetString("loglevel")),
		dishookBotMsg:  setDefaultString(fang.GetString("dishook.bot.message"), "Your HTTPS health monitoring result", true),
		dishookBotName: setDefaultString(fang.GetString("dishook.bot.name"), "HTTPS Doctor", true),
		dishookBotAva:  setDefaultString(fang.GetString("dishook.bot.avatar"), "https://www.a-trust.at/MediaProvider/2107/ssl.png", true),
		dishookUrl:     setDefaultString(fang.GetString("dishook.url"), "", true),
		location:       setLocation(fang.GetString("tz")),
		grace:          graceperiod,
	}
}

func splitCSV(s string) []string {
	strTrimmed := strings.TrimFunc(
		s,
		func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		},
	)
	comaToSpace := strings.NewReplacer(",", " ")
	strReplacedComa := comaToSpace.Replace(strTrimmed)
	rgxSpaceSplit := regexp.MustCompile(`\s+`)
	return rgxSpaceSplit.Split(strReplacedComa, -1)
}

func setLogLevel(l string) zerolog.Level {
	switch l {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.Disabled
	}
}

func setDefaultString(value, fallback string, trimSpace bool) string {
	if trimSpace {
		value = strings.TrimSpace(value)
	}
	if len(value) <= 0 {
		return fallback
	}
	return value
}

func setLocation(timezone string) *time.Location {
	loc, err := time.LoadLocation(strings.TrimSpace(timezone))
	if err != nil {
		loc = time.Local
	}
	return loc
}

func Get() Config {
	once.Do(func() {
		cfg = load()
	})
	return cfg
}
