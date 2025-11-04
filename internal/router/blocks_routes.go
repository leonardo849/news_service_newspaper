package router

import (
	"news_service/config"
	"news_service/internal/helper"
	"news_service/internal/logger"
	"news_service/internal/middleware"

	"github.com/gofiber/fiber/v2"
	constsSl "github.com/leonardo849/shared_library_news_paper/pkg/consts"
	middlewaresSl "github.com/leonardo849/shared_library_news_paper/pkg/middlewares"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func setupBlocksRoutes(blocksGroup fiber.Router, db *gorm.DB, rc *redis.Client) {
	blockController := helper.CreateBlockController(db, rc)
	blocksGroup.Post("/create/:news_id", middlewaresSl.VerifyJWT(config.Key), middlewaresSl.CheckRole([]string{constsSl.Journalist}),middleware.PermJournalist() ,blockController.CreateBlocks())
	logger.ZapLogger.Info("block's routes are working!")
}