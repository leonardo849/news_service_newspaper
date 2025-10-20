package helper

import (
	"news_service/internal/handler"
	"news_service/internal/repository"
	"news_service/internal/service"

	"gorm.io/gorm"
)

func createUserRepository(db *gorm.DB) *repository.UserRepository {
	return repository.CreateUserRepository(db)
}

func CreateUserService(db *gorm.DB) *service.UserService {
	return service.CreateNewUserService(createUserRepository(db))
}

func CreateUserController(db *gorm.DB) *handler.UserController {
	return handler.CreateNewUserController(CreateUserService(db))
}