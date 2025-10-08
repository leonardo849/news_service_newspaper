package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	News []NewsModel `json:"-" gorm:"many2many:authors_news;"`
	AuthId    string    `json:"auth_id"`
	Username  string    `gorm:"size:50;unique;not null" json:"username"`
	Role      string    `gorm:"default:'CUSTOMER';not null"`
	Bio       string    `gorm:"size:255;default:'hello'" json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
