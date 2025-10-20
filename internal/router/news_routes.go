package router

import (
	"news_service/config"
	"news_service/internal/helper"
	"news_service/internal/logger"

	"github.com/gofiber/fiber/v2"
	middlewaresSl "github.com/leonardo849/shared_library_news_paper/pkg/middlewares"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)



func setupNewsRoutes(newsGroup fiber.Router, db *gorm.DB, rc *redis.Client) {
	controller := helper.CreateNewsController(db, rc)
	newsGroup.Post("/create", middlewaresSl.VerifyJWT(config.Key), controller.CreateNews())
	newsGroup.Get("/one/:id", middlewaresSl.VerifyJWT(config.Key), controller.FindNewsById())
	newsGroup.Patch("/update/publish/:id", middlewaresSl.VerifyJWT(config.Key), controller.PublishNewsById())
	logger.ZapLogger.Info("news's routes are working!")
}