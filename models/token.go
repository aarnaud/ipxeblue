package models

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	Token         string        `gorm:"primaryKey" json:"token"`
	ComputerUUID  uuid.UUID     `gorm:"not null" json:"computer_uuid"`
	BootentryUUID uuid.UUID     `gorm:"not null" json:"bootentryfile_uuid"`
	Filename      string        `gorm:"not null" json:"filename"`
	BootentryFile BootentryFile `gorm:"ForeignKey:BootentryUUID,Filename;References:BootentryUUID,Name" json:"-"`
	Computer      Computer      `gorm:"ForeignKey:ComputerUUID;References:Uuid,Name" json:"-"`
	ExpireAt      time.Time     `json:"expire"`
}
