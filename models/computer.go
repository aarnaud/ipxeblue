package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"time"
)

type Computer struct {
	Uuid              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Mac               pgtype.Macaddr `gorm:"type:macaddr;index:idx_mac" json:"-"`
	IP                pgtype.Inet    `gorm:"type:inet;index:idx_ip" json:"-"`
	Asset             string         `json:"asset"`
	BuildArch         string         `json:"build_arch"`
	Hostname          string         `json:"hostname"`
	LastSeen          time.Time      `json:"last_seen"`
	Manufacturer      string         `json:"manufacturer"`
	Name              string         `json:"name"`
	Platform          string         `json:"platform"`
	Product           string         `json:"product"`
	Serial            string         `json:"serial"`
	Version           string         `json:"version"`
	Tags              []Tag          `gorm:"foreignkey:ComputerUUID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tags"`
	BootentryUUID     *uuid.UUID     `gorm:"type:uuid" json:"bootentry_uuid"`
	Bootentry         *Bootentry     `gorm:"foreignkey:bootentry_uuid;References:Uuid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;default:NULL"`
	LastIpxeaccountID string         `json:"last_ipxeaccount"`
	LastIpxeaccount   *Ipxeaccount   `gorm:"foreignkey:last_ipxeaccount_id;References:Username;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;default:NULL" json:"-"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

// MarshalJSON initializes nil slices and then marshals the bag to JSON
func (c Computer) MarshalJSON() ([]byte, error) {
	if c.Tags == nil {
		c.Tags = make([]Tag, 0)
	}

	type Alias Computer
	return json.Marshal(&struct {
		Mac string `json:"mac"`
		IP  string `json:"ip"`
		Alias
	}{
		Mac:   c.Mac.Addr.String(),
		IP:    c.IP.IPNet.IP.String(),
		Alias: (Alias)(c),
	})

}

func (c *Computer) UnmarshalJSON(data []byte) error {
	type Alias Computer
	aux := &struct {
		Mac string `json:"mac"`
		IP  string `json:"ip"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if err := c.Mac.DecodeText(nil, []byte(aux.Mac)); err != nil {
		return err
	}
	if err := c.IP.DecodeText(nil, []byte(aux.IP)); err != nil {
		return err
	}
	return nil
}
