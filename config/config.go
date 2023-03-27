package config

import (
	"log"

	"github.com/spf13/viper"
)

type Settings struct {
	EncKey string `mapstructure:"ENC_KEY"`
}

var Env *Settings

func New() *Settings {
	var cfg Settings

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("No env file, using environment variables.", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal("Error trying to unmarshal configuration", err)
	}
	
	Env = &cfg
	return Env
}
