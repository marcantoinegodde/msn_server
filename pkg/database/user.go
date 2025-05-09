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
	WebauthnID          []byte `gorm:"uniqueIndex"`
	WebauthnCredentials []Credential
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
	credentials := make([]webauthn.Credential, len(u.WebauthnCredentials))
	for i, cred := range u.WebauthnCredentials {
		credentials[i] = webauthn.Credential{
			ID:              cred.KeyID,
			PublicKey:       cred.PublicKey,
			AttestationType: cred.AttestationType,
			Transport:       cred.Transport,
			Flags: webauthn.CredentialFlags{
				UserPresent:    cred.UserPresent,
				UserVerified:   cred.UserVerified,
				BackupEligible: cred.BackupEligible,
				BackupState:    cred.BackupState,
			},
			Authenticator: webauthn.Authenticator{
				AAGUID:       cred.AAGUID,
				SignCount:    cred.SignCount,
				CloneWarning: cred.CloneWarning,
				Attachment:   cred.Attachment,
			},
			Attestation: webauthn.CredentialAttestation{
				ClientDataJSON:     cred.ClientDataJSON,
				ClientDataHash:     cred.ClientDataHash,
				AuthenticatorData:  cred.AuthenticatorData,
				PublicKeyAlgorithm: cred.PublicKeyAlgorithm,
				Object:             cred.Object,
			},
		}
	}
	return credentials
}
