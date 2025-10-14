package helper

import (
	"news_service/internal/repository"

	"gorm.io/gorm"
)

func createUserRepository(db *gorm.DB) *repository.UserRepository {
	return repository.CreateUserRepository(db)
}