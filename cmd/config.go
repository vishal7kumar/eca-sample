package cmd

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var config *viper.Viper //TODO: should be written in a wrapper ?

func initConfig() {
	config = viper.New()  //TODO: should need to add options ?
	config.AutomaticEnv() // read value ENV variable
	stage := config.GetString("STAGE")
	config.SetConfigName(stage) // TODO: config file should be based on ENV
	config.SetConfigType("json")
	config.AddConfigPath(".")
	config.AddConfigPath("./config/") // config file path

	err := config.ReadInConfig()
	if err != nil {
		log.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
}

func GetConfig() *viper.Viper {
	if config == nil {
		initConfig()
	}

	return config
}
