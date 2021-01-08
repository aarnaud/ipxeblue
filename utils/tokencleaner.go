package utils

import (
	"fmt"
	"github.com/aarnaud/ipxeblue/models"
	"gorm.io/gorm"
	"time"
)

func TokenCleaner(db *gorm.DB) {
	token := models.Token{}
	for {
		result := db.Where("expire_at < NOW()").Delete(&token)
		if result.RowsAffected > 0 {
			fmt.Printf("%d tokens was expired and deleted", result.RowsAffected)
		}
		time.Sleep(time.Second)
	}
}
