package models

import (
	"time"

	"gorm.io/gorm"
)

type StorageProvider struct {
	CreatedAt      time.Time         `json:"created_at,omitempty"`
	UpdatedAt      time.Time         `json:"updated_at,omitempty"`
	DeletedAt      gorm.DeletedAt    `gorm:"index" json:"deleted_at,omitempty"`
	Address        string            `gorm:"primaryKey" json:"address"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	ReputationLink string            `json:"reputation_link"`
	AvatarURL      string            `json:"avatar_url"`
	Proposals      []BigFileProposal `json:"proposals" gorm:"foreignKey:StorageProviderAddress;references:Address"`
	ContactInfo    string            `json:"contact_info"`
}
