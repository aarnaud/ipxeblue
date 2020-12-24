package models

import "github.com/google/uuid"

type Tag struct {
	Key          string    `gorm:"primaryKey;index" json:"key"`
	Value        string    `json:"value"`
	ComputerUUID uuid.UUID `gorm:"primaryKey;index" json:"-"`
}
