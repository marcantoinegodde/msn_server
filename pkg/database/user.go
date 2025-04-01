package database

import (
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email               string `gorm:"uniqueIndex"`
	Salt                string
	Password            string
	WebauthnID          []byte                `gorm:"uniqueIndex"`
	WebauthnCredentials []webauthn.Credential `gorm:"serializer:json"`
	FirstName           string
	LastName            string
	Country             string
	State               string
	City                string
	DisplayName         string
	Status              Status  `gorm:"default:FLN"`
	DataVersion         uint32  `gorm:"default:0"`
	Gtc                 Gtc     `gorm:"default:A"`
	Blp                 Blp     `gorm:"default:AL"`
	ForwardList         []*User `gorm:"many2many:forward_list"`
	AllowList           []*User `gorm:"many2many:allow_list"`
	BlockList           []*User `gorm:"many2many:block_list"`
	ReverseList         []*User `gorm:"many2many:reverse_list"`
}

func (u *User) WebAuthnID() []byte {
	return u.WebauthnID
}

func (u *User) WebAuthnName() string {
	return u.Email
}

func (u *User) WebAuthnDisplayName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.WebauthnCredentials
}
