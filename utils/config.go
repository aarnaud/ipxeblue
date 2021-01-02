package utils

import "github.com/spf13/viper"

type MinioConfig struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
	Secure     bool
}

type Config struct {
	Port          int
	EnableAPIAuth bool
	MinioConfig   MinioConfig
	BaseURL       string
}

func GetConfig() *Config {
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv() // To get the value from the config file using key// viper package read .env

	config := Config{
		Port:          8080,
		EnableAPIAuth: true,
		MinioConfig: MinioConfig{
			Endpoint:   "127.0.0.1:9000",
			BucketName: "ipxeblue",
		},
		BaseURL: "http://127.0.0.1:8080",
	}

	if p := viper.GetInt("PORT"); p != 0 {
		config.Port = p
	}

	if value := viper.GetString("MINIO_ENDPOINT"); value != "" {
		config.MinioConfig.Endpoint = value
	}

	config.MinioConfig.AccessKey = viper.GetString("MINIO_ACCESS_KEY")
	config.MinioConfig.SecretKey = viper.GetString("MINIO_SECRET_KEY")
	config.MinioConfig.Secure = viper.GetBool("MINIO_SECURE")

	if value := viper.GetString("MINIO_BUCKETNAME"); value != "" {
		config.MinioConfig.BucketName = value
	}

	if viper.IsSet("ENABLE_API_AUTH") {
		config.EnableAPIAuth = viper.GetBool("ENABLE_API_AUTH")
	}

	if value := viper.GetString("BASE_URL"); value != "" {
		config.BaseURL = value
	}

	return &config
}
