package models

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	Token         string    `gorm:"primaryKey" json:"token"`
	Computer      Computer  `gorm:"ForeignKey:ComputerUUID;References:Uuid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ComputerUUID  uuid.UUID `gorm:"not null" json:"computer_uuid"`
	Bootentry     Bootentry `gorm:"ForeignKey:BootentryUUID;References:Uuid,Name;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	BootentryUUID uuid.UUID `gorm:"not null" json:"bootentry_uuid"`
	// BootentryFile can be null if we generate DownloadBaseURL
	BootentryFile *BootentryFile `gorm:"ForeignKey:BootentryUUID,Filename;References:BootentryUUID,Name;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Filename      *string        `json:"filename"`
	ExpireAt      time.Time      `gorm:"index:idx_expire_at" json:"expire"`
}
