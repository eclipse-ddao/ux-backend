package models

import (
	"time"

	"gorm.io/gorm"
)

type BigFile struct {
	ID                 uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt          time.Time         `json:"created_at,omitempty"`
	UpdatedAt          time.Time         `json:"updated_at,omitempty"`
	DeletedAt          gorm.DeletedAt    `gorm:"index" json:"deleted_at"`
	Duration           uint              `json:"duration"`
	Expiry             time.Time         `json:"expiry"`
	SizeInGb           uint              `json:"size_in_gb"`
	BaseBounty         uint              `json:"base_bounty"`
	FileType           uint              `json:"file_type"` // 1 = Link, 2 = Hardware Delivery, 3 = Filecoin CID
	Name               string            `json:"name"`
	Description        string            `json:"description"`
	DaoContractAddress string            `json:"dao_contract_address"`
	UploadedBy         string            `json:"uploaded_by"`
	Proposals          []BigFileProposal `json:"proposals" gorm:"foreignKey:BFileID"`
	SelectedProposalID uint              `json:"selected_proposal_id"`
	Status             uint              `json:"status"` // 1 = Open, 2 = SP Selected, 3 = Close (SP Accepted)
}
