package models

import (
	"fmt"
	"github.com/aarnaud/ipxeblue/utils/helpers"
	"github.com/google/uuid"
	"time"
)

type Bootentry struct {
	Uuid        uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string          `gorm:"uniqueIndex:idx_name" json:"name"`
	Description string          `json:"description"`
	Persistent  *bool           `gorm:"not null;default:FALSE" json:"persistent"`
	IpxeScript  string          `json:"ipxe_script"`
	Files       []BootentryFile `gorm:"foreignkey:bootentry_uuid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"files"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (b *Bootentry) GetFile(filename string) *BootentryFile {
	for _, file := range b.Files {
		if file.Name == filename {
			file.Bootentry = *b
			return &file
		}
	}
	return nil
}

func (b *Bootentry) GetDownloadBasePath() (string, *Token) {
	token := Token{
		Token:         helpers.RandomString(15),
		Bootentry:     *b,
		BootentryFile: nil,
		// TODO: expose token duration in configuration
		ExpireAt: time.Now().Add(time.Minute * 10),
	}
	return fmt.Sprintf("/files/token/%s/%s/", token.Token, b.Uuid), &token
}
