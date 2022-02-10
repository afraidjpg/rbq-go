package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var Cfg *viper.Viper

func init() {
	loadConfig()
}

func loadConfig() {
	Cfg = viper.New()
	Cfg.SetConfigName("config")
	Cfg.SetConfigType("yaml")
	Cfg.AddConfigPath(".")
	Cfg.AddConfigPath("./../../")

	err := Cfg.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}
