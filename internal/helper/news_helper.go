package helper

import (
	"news_service/internal/handler"
	"news_service/internal/repository"
	"news_service/internal/service"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func createNewsRedisRepository(rc *redis.Client) *repository.NewsRedisRepository {
	return  repository.CreateNewsRedisRepository(rc)
}

func createNewsRepository(db *gorm.DB) *repository.NewsRepository {
	return repository.CreateNewsRepository(db)
}

func CreateNewsService(db *gorm.DB, rc *redis.Client) *service.NewsService {
	return service.CreateNewsService(createNewsRepository(db), createUserRepository(db), CreateUnitOfWork(db), createNewsRedisRepository(rc))
}

func CreateNewsController(db *gorm.DB, rc *redis.Client) *handler.NewsController {
	return handler.CreateNewsController(CreateNewsService(db, rc))
}