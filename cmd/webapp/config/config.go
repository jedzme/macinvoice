package config

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"strings"
)

// AppConfig is for the configuration values.
var (
	AppConfig *viper.Viper
)

//
func init() {

	AppConfig = viper.New()

	//AppConfig.AddConfigPath("../../")
	//AppConfig.SetConfigName("settings")
	//AppConfig.SetConfigType("yaml")
	AppConfig.SetConfigFile("../../settings.yaml")
	AppConfig.ReadInConfig()

	setLog()
}

func setLog() {

	l := strings.ToUpper(AppConfig.GetString("LOG_LEVEL"))

	switch l {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "TRACE":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "FATAL":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "PANIC":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	fmt.Println("{ Log Level is set to: " + l + "}")

}
