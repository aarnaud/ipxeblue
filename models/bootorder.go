package models

import (
	"github.com/google/uuid"
	"time"
)

type Bootorder struct {
	Order         int        `gorm:"not null;" json:"-"`
	ComputerUuid  uuid.UUID  `gorm:"primaryKey;not null;"`
	Computer      *Computer  `gorm:"foreignkey:computer_uuid;References:Uuid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	BootentryUuid uuid.UUID  `gorm:"primaryKey;not null;"`
	Bootentry     *Bootentry `gorm:"foreignkey:bootentry_uuid;References:Uuid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreatedAt     time.Time  `gorm:"autoCreateTime;not null;default:current_timestamp" json:"-"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime;not null;default:current_timestamp" json:"-"`
}
