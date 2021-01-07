package utils

import (
	"fmt"
	"github.com/aarnaud/ipxeblue/models"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Database() *gorm.DB {
	var err error
	var db *gorm.DB
	var databaseUrl string

	viperDBUrl := viper.GetString("DATABASE_URL")

	if viperDBUrl == "" {
		DBuser := viper.GetString("DB_USER")
		DBpassword := viper.GetString("DB_PASSWORD")
		DBname := viper.GetString("DB_NAME")
		DBhost := viper.GetString("DB_HOST")
		DBport := viper.GetString("DB_PORT")
		DBsslmode := viper.GetString("DB_SSLMODE")
		databaseUrl = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			DBuser, DBpassword, DBhost, DBport, DBname, DBsslmode)
	} else {
		databaseUrl = viperDBUrl
	}

	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: databaseUrl,
	}), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic("Failed to connect to database!")
	}

	err = db.AutoMigrate(&models.Computer{}, &models.Tag{}, &models.Ipxeaccount{}, &models.Bootentry{},
		&models.BootentryFile{}, &models.Token{})

	if err != nil {
		fmt.Println(err)
		panic("Failed to migrate database!")
	}
	return db
}
