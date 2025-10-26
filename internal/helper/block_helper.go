package helper

import (
	"news_service/internal/handler"
	"news_service/internal/repository"
	"news_service/internal/service"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func createBlockRepository(db *gorm.DB) *repository.BlockRepository{
	return repository.CreateBlockRepository(db)
}

func CreateBlockService(db *gorm.DB, rc *redis.Client) *service.BlockService {
	return service.CreateBlockService(CreateNewsService(db, rc), CreateUnitOfWork(db))
}

func CreateBlockController(db *gorm.DB, rc *redis.Client) *handler.BlockController {
	return handler.CreateBlockController(CreateBlockService(db, rc))
}