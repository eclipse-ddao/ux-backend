package models

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID                 uint           `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	CreatedAt          time.Time      `json:"created_at,omitempty"`
	UpdatedAt          time.Time      `json:"updated_at,omitempty"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name               string         `json:"name"`
	UploadedBy         string         `json:"uploaded_by"`
	ImageURL           string         `json:"image_url"`
	DaoContractAddress string         `json:"dao_contract_address"`
	Cid                string         `json:"cid"`
	FileType           string         `json:"file_type"`
}
