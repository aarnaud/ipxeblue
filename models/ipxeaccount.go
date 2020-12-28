package models

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type Ipxeaccount struct {
	Username  string    `gorm:"primarykey" json:"username"`
	Password  string    `json:"password,omitempty"`
	LastLogin time.Time `json:"last_login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HashPassword : Encrypt user password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// BeforeSave : hook before a user is saved
func (o *Ipxeaccount) BeforeSave(tx *gorm.DB) (err error) {
	if o.Password != "" {
		hash, err := HashPassword(o.Password)
		if err != nil {
			return nil
		}
		o.Password = hash
	}

	return
}

// MarshalJSON initializes nil slices and then marshals the bag to JSON
func (o Ipxeaccount) MarshalJSON() ([]byte, error) {
	type Alias Ipxeaccount
	alias := (Alias)(o)
	// empty password
	alias.Password = ""
	// Add Id for react-admin
	return json.Marshal(&struct {
		Id string `json:"id"`
		Alias
	}{
		Id: alias.Username,
		Alias: alias,
	})

}

func (c *Ipxeaccount) UnmarshalJSON(data []byte) error {
	type Alias Ipxeaccount
	aux := &struct {
		PasswordConfirmation string `json:"password_confirmation"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Password != "" && aux.Password != aux.PasswordConfirmation {
		return fmt.Errorf("Password missmatch")
	}
	return nil
}
