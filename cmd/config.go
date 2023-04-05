package cmd

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

var config *viper.Viper

func initConfig() {
	config = viper.New()
	config.AutomaticEnv() // read value ENV variable
	stage := config.GetString("STAGE")
	config.SetConfigName(stage)
	config.SetConfigType("json")
	config.AddConfigPath(".")
	config.AddConfigPath("./config/")  // config file path
	config.AddConfigPath("../config/") // config file path

	err := config.ReadInConfig()
	if err != nil {
		log.Println("fatal error config file: \n", err)
		os.Exit(1)
	}
}

func GetConfig() *viper.Viper {
	if config == nil {
		initConfig()
	}

	return config
}
