package helper

import (
	"news_service/internal/repository"

	"gorm.io/gorm"
)

func CreateBlockRepository(db *gorm.DB) *repository.BlockRepository{
	return repository.CreateBlockRepository(db)
}