package utils

import (
	"fmt"
	"github.com/aarnaud/ipxeblue/models"
	zlog "github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
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

	dbLogger := logger.New(
		log.New(os.Stdout, "", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Warn,            // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,                  // Disable color
		},
	)

	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: databaseUrl,
	}), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		zlog.Panic().Err(err).Msg("failed to connect to database!")
	}

	err = db.AutoMigrate(&models.Computer{}, &models.Tag{}, &models.Ipxeaccount{}, &models.Bootentry{},
		&models.BootentryFile{}, &models.Bootorder{}, &models.Token{})
	if err != nil {
		zlog.Panic().Err(err).Msg("failed to automigrate database!")
	}

	// Custom migration database schema
	// remove bootentry_uuid in computer since we have bootorder
	computerRef := &models.Computer{}
	if db.Migrator().HasColumn(computerRef, "bootentry_uuid") {
		if err = db.Migrator().DropColumn(computerRef, "bootentry_uuid"); err != nil {
			zlog.Panic().Err(err).Msg("failed to drop column bootentry_uuid")
		}
	}

	return db
}
