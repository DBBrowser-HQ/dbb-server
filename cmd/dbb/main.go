package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("Error while init configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env file: %s", err.Error())
	}

	logrus.Info(viper.Get("bind_addr"))
	logrus.Info(os.Getenv("PASSWORD_SALT"))
}

func initConfig() error {
	viper.AddConfigPath("cmd/dbb/config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
