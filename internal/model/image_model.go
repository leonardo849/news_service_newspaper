package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImageModel struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	URL string `json:"url" gorm:"not null"`
	BlockID uuid.UUID `json:"block_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (i *ImageModel) BeforeCreate(tx *gorm.DB) (err error) {
	i.ID = uuid.New()
	return
}