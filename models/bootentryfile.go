package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type BootentryFile struct {
	Name          string     `gorm:"primaryKey;index" json:"name"`
	Protected     *bool      `gorm:"not null;default:FALSE" json:"protected"`
	Templatized   *bool      `gorm:"not null;default:FALSE" json:"templatized"`
	BootentryUUID uuid.UUID  `gorm:"type:uuid;primaryKey;index" json:"-"`
	Bootentry     *Bootentry `gorm:"foreignkey:bootentry_uuid;References:Uuid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (b *BootentryFile) GetFileStorePath() string {
	return fmt.Sprintf("%s/files/%s", b.BootentryUUID.String(), b.Name)
}

func (b *BootentryFile) GetAPIDownloadPath() string {
	return fmt.Sprintf("/api/v1/bootentries/%s/files/%s", b.BootentryUUID.String(), b.Name)
}

func (b *BootentryFile) GetDownloadPath() (string, *Token) {
	basePath, token := b.Bootentry.GetDownloadBasePath()
	if *b.Protected {
		token.BootentryFile = b
		return fmt.Sprintf("%s%s", basePath, b.Name), token
	}
	return fmt.Sprintf("/files/public/%s/%s", b.BootentryUUID.String(), b.Name), nil
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

func (o *BootentryFile) UnmarshalJSON(data []byte) error {

	type Alias BootentryFile
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	falseRef := false
	if o.Protected == nil {
		o.Protected = &falseRef
	}
	if o.Templatized == nil {
		o.Protected = &falseRef
	}
	return nil
}
