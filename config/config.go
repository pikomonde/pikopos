package config

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	MySQLPikopos string
	JWTSecret    string
	BaseURL      string

	OAuth map[string]struct {
		ClientID     string
		ClientSecret string
		Scopes       []string
	}
}

// TODO: remove it, return the variable to app.go instead. These struct
// should only be used in clients, not be used in services or deliveries.
var C Config

// TODO: change from viper to other lighter library

// Init is used to initialize new config
func Init() {
	appPath, err := os.Executable()
	if err != nil {
		log.WithFields(log.Fields{
			"appPath": appPath,
		}).Panicln("[Config]: Failed getting executable app location", err.Error())
	}
	appPathSplit := strings.Split(appPath, "/")
	if len(appPathSplit) < 1 {
		log.WithFields(log.Fields{
			"appPath":      appPath,
			"appPathSplit": appPathSplit,
		}).Panicln("[Config]: Invalid xecutable app location")
	}
	appPathDir := strings.Join(appPathSplit[:len(appPathSplit)-1], "/")

	if os.Getenv("env") == "PROD" {
		viper.SetConfigFile(appPathDir + "/config/production.yaml")
	} else {
		viper.SetConfigFile(appPathDir + "/config/local.yaml")
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Panicln("[Config]: Failed initialize config", err.Error())
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		log.Panicln("[Config]: Failed unmarshal config", err.Error())
	}
}
