package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Port string
}


func SetConfig() *Config{

	if os.Getenv("config") != "" {
		viper.SetConfigFile(fmt.Sprintf("./config/%s",os.Getenv("config")))
	} else {
		viper.SetConfigFile("./config/dev.config.yaml")
	}

	viper.SetDefault("port", ":3000")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to Read in Config File: %s", err.Error())
	}
	return &Config{
		Port: viper.GetString("port"),
	}
}