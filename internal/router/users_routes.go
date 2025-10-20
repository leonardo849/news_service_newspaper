package router

import (
	"news_service/config"
	"news_service/internal/helper"
	"news_service/internal/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/leonardo849/shared_library_news_paper/pkg/middlewares"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func setupUserRoutes(usersGroup fiber.Router, db *gorm.DB, rc *redis.Client) {
	userController := helper.CreateUserController(db)
	usersGroup.Patch("/update/bio", middlewares.VerifyJWT(config.Key), userController.UpdateBio())
	logger.ZapLogger.Info("user's routes are ready")
}
