package utils

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net/url"
)

type MinioConfig struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
	Secure     bool
}

type Config struct {
	Port                 int
	EnableAPIAuth        bool
	MinioConfig          MinioConfig
	BaseURL              *url.URL
	GrubSupportEnabled   bool
	TFTPEnabled          bool
	DefaultBootentryName string
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
		GrubSupportEnabled: false,
		TFTPEnabled:        false,
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

	BaseURL := "http://127.0.0.1:8080"
	if value := viper.GetString("BASE_URL"); value != "" {
		BaseURL = value
	}

	u, err := url.Parse(BaseURL)
	if err != nil {
		log.Panic().Err(err).Msg("failed to parse BASE_URL")
	}
	config.BaseURL = u

	config.GrubSupportEnabled = viper.GetBool("GRUB_SUPPORT_ENABLED")
	config.TFTPEnabled = viper.GetBool("TFTP_ENABLED")
	config.DefaultBootentryName = viper.GetString("DEFAULT_BOOTENTRY_NAME")

	return &config
}
