package helper

import (
	"news_service/internal/repository"

	"gorm.io/gorm"
)

func CreateImageRepository(db *gorm.DB) *repository.ImageRepository {
	return repository.CreateImageRepository(db)
}