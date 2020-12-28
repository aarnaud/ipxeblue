package config

import "github.com/spf13/viper"

type Config struct {
	Port int
	EnableAPIAuth bool
}


func GetConfig() Config {
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv() // To get the value from the config file using key// viper package read .env

	config := Config{
		Port: 8080,
		EnableAPIAuth: true,
	}

	if p := viper.GetInt("PORT"); p != 0 {
		config.Port = p
	}

	if viper.IsSet("ENABLE_API_AUTH") {
		config.EnableAPIAuth = viper.GetBool("ENABLE_API_AUTH")
	}

	return config
}