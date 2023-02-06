package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type User struct {
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Address   string         `gorm:"primaryKey" json:"address"`
	Username  string         `json:"username"`
	AvatarURL string         `json:"avatar_url"`
	Daos      []*Dao         `gorm:"many2many:user_daos;" json:"daos"`
	BFiles    []BigFile      `gorm:"foreignKey:UploadedBy;references:Address" json:"b_files"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("address", strings.ToUpper(u.Address))
	return nil
}

func (u *User) AfterFind(tx *gorm.DB) (err error) {
	u.Address = strings.ToUpper(u.Address)
	return nil
}
