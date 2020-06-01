package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	MySQLPikopos string
	JWTSecret    string
}

var C Config

// TODO: change from viper to other lighter library

// Init is used to initialize new config
func Init() {
	if os.Getenv("env") == "PROD" {
		viper.SetConfigFile("config/production.yaml")
	} else {
		viper.SetConfigFile("config/local.yaml")
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Panicln("[Config]: Failed initialize config", err.Error())
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		log.Panicln("[Config]: Failed unmarshal config", err.Error())
	}
}
