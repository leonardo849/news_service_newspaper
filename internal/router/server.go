package router

import (
	"news_service/config"
	_ "news_service/internal/dto"
	"news_service/internal/logger"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	constsSl "github.com/leonardo849/shared_library_news_paper/pkg/consts"
	middlewaresSl "github.com/leonardo849/shared_library_news_paper/pkg/middlewares"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupApp(db *gorm.DB, rc *redis.Client) *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	
	logger.ZapLogger.Info("cors is ready")
	app.Use(middlewaresSl.LogRequestsMiddleware())
	
	// @Summary Hello
	// @Description welcome message
	// @Accept json
	// @Produce json
	// @Sucess 200 {object} dto.MessageDTO
	// @Router / [get]
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON(fiber.Map{"message": "what's up?"})
	})

	
	app.Get("/swagger/*",  middlewaresSl.VerifyJWT(config.Key), middlewaresSl.CheckRole([]string{constsSl.Developer, constsSl.Ceo}) ,swagger.HandlerDefault)
	app.Get("/metrics", middlewaresSl.VerifyJWT(config.Key), middlewaresSl.CheckRole([]string{constsSl.Developer, constsSl.Ceo}) ,adaptor.HTTPHandler(promhttp.Handler()))
	logger.ZapLogger.Info("swagger and prometheus are ready")

	newsGroup := app.Group("/news")
	usersGroup := app.Group("/users")
	blocksGroup := app.Group("/blocks")
	setupNewsRoutes(newsGroup, db, rc)
	setupUserRoutes(usersGroup, db, rc)
	setupBlocksRoutes(blocksGroup, db, rc)
	logger.ZapLogger.Info("app is running!")
	return  app
}

func RunServer(db *gorm.DB, rc *redis.Client) error {
	app := SetupApp(db, rc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	return app.Listen(":" + port)
}