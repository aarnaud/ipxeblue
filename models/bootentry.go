package models

import (
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

func (b *Bootentry) GetDownloadPath(filename string) string {
	for _, file := range b.Files {
		if file.Name == filename {
			return file.GetDownloadPath()
		}
	}
	return ""
}
