package models

import (
	"time"

	"gorm.io/gorm"
)

type BigFileProposal struct {
	ID                     uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt              time.Time      `json:"created_at,omitempty"`
	UpdatedAt              time.Time      `json:"updated_at,omitempty"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	RequestedBounty        uint           `json:"requested_bounty"`
	DealCID                string         `json:"deal_cid"`
	IsSelected             bool           `json:"is_selected"`
	StorageProviderAddress string         `json:"storage_provider_address"`
	BFileID                uint           `json:"b_file_id"`
	BFile                  *BigFile       `json:"b_file"`
}
