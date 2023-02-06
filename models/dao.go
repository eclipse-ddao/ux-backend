package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Dao struct {
	ContractAddress string         `gorm:"primaryKey" json:"contract_address,omitempty"`
	CreatedAt       time.Time      `json:"created_at,omitempty"`
	UpdatedAt       time.Time      `json:"updated_at,omitempty"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	AvatarURL       string         `json:"avatar_url"`
	Members         []*User        `gorm:"many2many:user_daos;" json:"members"`
	Files           []File         `gorm:"foreignKey:DaoContractAddress;references:ContractAddress" json:"files"`
	BFiles          []BigFile      `gorm:"foreignKey:DaoContractAddress;references:ContractAddress" json:"big_files"`
}

func (u *Dao) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("contract_address", strings.ToUpper(u.ContractAddress))
	return nil
}

func (u *Dao) AfterFind(tx *gorm.DB) (err error) {
	u.ContractAddress = strings.ToUpper(u.ContractAddress)
	return nil
}
