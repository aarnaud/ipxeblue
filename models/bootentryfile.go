package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type BootentryFile struct {
	Name          string    `gorm:"primaryKey;index" json:"name"`
	Protected     *bool     `json:"protected"`
	BootentryUUID uuid.UUID `gorm:"type:uuid;primaryKey;index" json:"-"`
}

func (b *BootentryFile) GetFileStorePath() string {
	return fmt.Sprintf("%s/files/%s", b.BootentryUUID.String(), b.Name)
}

func (b *BootentryFile) GetAPIDownloadPath() string {
	return fmt.Sprintf("/api/v1/bootentries/%s/files/%s", b.BootentryUUID.String(), b.Name)
}

func (b *BootentryFile) GetDownloadPath() string {
	return fmt.Sprintf("/files/%s/%s", b.BootentryUUID.String(), b.Name)
}

// MarshalJSON initializes nil slices and then marshals the bag to JSON
func (o BootentryFile) MarshalJSON() ([]byte, error) {
	type ReactFile struct {
		Src   string `json:"src"`
		Title string `json:"title"`
	}
	type Alias BootentryFile
	alias := (Alias)(o)

	// Add Id for react-admin
	return json.Marshal(&struct {
		File ReactFile `json:"file"`
		Alias
	}{
		File: ReactFile{
			Src:   o.GetAPIDownloadPath(),
			Title: alias.Name,
		},
		Alias: alias,
	})
}
