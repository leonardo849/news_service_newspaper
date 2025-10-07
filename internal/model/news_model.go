package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NewsModel struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Authors []UserModel `json:"authors" gorm:"many2many:authors_news;"`
	Title string `json:"title" gorm:"size:50;not null"`
	Subtitle string `json:"subtitle" gorm:"size:100;not null"`
	Topic string `json:"topic" gorm:"size:100;not null"`
	Status string `json:"status" gorm:"default:'SKETCH';not null"`
	Published_at *time.Time `json:"published_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (n *NewsModel) BeforeCreate(tx *gorm.DB) (err error) {
	n.ID = uuid.New()
	return
}