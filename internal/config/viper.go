package config

import (
	"github.com/spf13/viper"
	"github.com/zekroTJA/echo/internal/verbosity"
)

func InitViper() {
	viper.SetDefault(KeyAddr, ":80")
	viper.SetDefault(KeyVerbosity, verbosity.Normal)
	viper.SetDefault(KeyDebug, false)

	viper.SetEnvPrefix("ECHO")
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/echo/")
	viper.AddConfigPath("$HOME/.echo")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
}
