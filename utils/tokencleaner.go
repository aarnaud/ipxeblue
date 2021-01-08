package utils

import (
	"github.com/aarnaud/ipxeblue/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

func TokenCleaner(db *gorm.DB) {
	token := models.Token{}
	for {
		result := db.Where("expire_at < NOW()").Delete(&token)
		if result.RowsAffected > 0 {
			log.Info().Msgf("%d tokens was expired and deleted", result.RowsAffected)
		}
		time.Sleep(time.Second)
	}
}
