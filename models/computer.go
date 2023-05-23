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
	Tags              []*Tag         `gorm:"foreignkey:computer_uuid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tags"`
	Bootorder         []*Bootorder   `gorm:"foreignKey:computer_uuid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	LastIpxeaccountID *string        `json:"last_ipxeaccount"`
	LastIpxeaccount   *Ipxeaccount   `gorm:"foreignkey:last_ipxeaccount_id;References:Username;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;default:NULL" json:"-"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

var ComputerSearchFields = []string{
	"name",
	"hostname",
	"mac",
	"ip",
	"serial",
}

// MarshalJSON initializes nil slices and then marshals the bag to JSON
func (c Computer) MarshalJSON() ([]byte, error) {
	if c.Tags == nil {
		c.Tags = make([]*Tag, 0)
	}

	type Alias Computer

	// List of Bootentry ID for json output
	bootorder := make([]*Bootentry, len(c.Bootorder))
	for i, bootentry := range c.Bootorder {
		bootorder[i] = bootentry.Bootentry
	}

	return json.Marshal(&struct {
		Mac       string       `json:"mac"`
		IP        string       `json:"ip"`
		Bootorder []*Bootentry `json:"bootorder"`
		Alias
	}{
		Mac:       c.Mac.Addr.String(),
		IP:        c.IP.IPNet.IP.String(),
		Bootorder: bootorder,
		Alias:     (Alias)(c),
	})

}

func (c *Computer) UnmarshalJSON(data []byte) error {
	type Alias Computer
	aux := &struct {
		Mac       string       `json:"mac"`
		IP        string       `json:"ip"`
		Bootorder []*Bootentry `json:"bootorder"`
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
	// build the bootorder list
	c.Bootorder = make([]*Bootorder, len(aux.Bootorder))
	for i, bootentry := range aux.Bootorder {
		c.Bootorder[i] = &Bootorder{
			Order:         i,
			BootentryUuid: bootentry.Uuid,
			ComputerUuid:  c.Uuid,
		}
	}
	return nil
}
