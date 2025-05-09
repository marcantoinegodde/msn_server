package database

import (
	"github.com/go-webauthn/webauthn/protocol"
	"gorm.io/gorm"
)

type Credential struct {
	gorm.Model

	UserID             uint
	Name               string `gorm:"not null;default:''"`
	KeyID              []byte `gorm:"uniqueIndex;not null"`
	PublicKey          []byte `gorm:"not null"`
	AttestationType    string
	Transport          []protocol.AuthenticatorTransport `gorm:"serializer:json"`
	UserPresent        bool                              `gorm:"not null;default:false"`
	UserVerified       bool                              `gorm:"not null;default:false"`
	BackupEligible     bool                              `gorm:"not null;default:false"`
	BackupState        bool                              `gorm:"not null;default:false"`
	AAGUID             []byte                            `gorm:"not null"`
	SignCount          uint32                            `gorm:"not null;default:0"`
	CloneWarning       bool                              `gorm:"not null;default:false"`
	Attachment         protocol.AuthenticatorAttachment  `gorm:"not null"`
	ClientDataJSON     []byte                            `gorm:"not null"`
	ClientDataHash     []byte                            `gorm:"not null"`
	AuthenticatorData  []byte                            `gorm:"not null"`
	PublicKeyAlgorithm int64                             `gorm:"not null"`
	Object             []byte                            `gorm:"not null"`
}
