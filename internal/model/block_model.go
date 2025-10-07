package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BlockModel struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Content string `gorm:"size:1000;not null"`
	Images []ImageModel `json:"images" gorm:"foreignKey:BlockID;references:ID;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}	

func (b *BlockModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}
