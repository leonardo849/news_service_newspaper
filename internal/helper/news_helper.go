package helper

import (
	"news_service/internal/handler"
	"news_service/internal/repository"
	"news_service/internal/service"

	"gorm.io/gorm"
)

func createNewsRepository(db *gorm.DB) *repository.NewsRepository {
	return repository.CreateNewsRepository(db)
}

func CreateNewsService(db *gorm.DB) *service.NewsService {
	return service.CreateNewsService(createNewsRepository(db), createUserRepository(db), CreateUnitOfWork(db))
}

func CreateNewsController(db *gorm.DB) *handler.NewsController {
	return handler.CreateNewsController(CreateNewsService(db))
}