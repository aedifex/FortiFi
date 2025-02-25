package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Port 			string
	DB_USER 		string
	DB_PASS 		string
	DB_URL 			string
	DB_NAME 		string
	SIGNING_KEY 	string
	FcmKeyPath		string
	CORS_ORIGIN 	string
	OpenAIKey		string
}


func SetConfig() *Config{

	if os.Getenv("config") != "" {
		viper.SetConfigFile(fmt.Sprintf("./config/%s.config.yaml",os.Getenv("config")))
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
		DB_USER: viper.GetString("db_user"),
		DB_PASS: viper.GetString("db_pass"),
		DB_URL: viper.GetString("db_url"),
		DB_NAME: viper.GetString("db_name"),
		SIGNING_KEY: viper.GetString("signing_key"),
		FcmKeyPath: viper.GetString("fcm_key_path"),
		CORS_ORIGIN: viper.GetString("cors_origin"),
		OpenAIKey: viper.GetString("openai_key"),
	}
}