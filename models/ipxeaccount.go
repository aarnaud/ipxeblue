package models

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Ipxeaccount struct {
	Username  string    `gorm:"primarykey" json:"username"`
	Password  string    `json:"password,omitempty"`
	IsAdmin   *bool     `gorm:"not null;default:FALSE" json:"is_admin"`
	LastLogin time.Time `json:"last_login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HashPassword : Encrypt user password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
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
		Id:    alias.Username,
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
	if aux.Password != "" {
		hash, err := HashPassword(aux.Password)
		if err != nil {
			return err
		}
		aux.Password = hash
	}
	falseRef := false
	if c.IsAdmin == nil {
		c.IsAdmin = &falseRef
	}
	return nil
}
