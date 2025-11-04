package middleware

import (
	"news_service/internal/helper"
	"news_service/internal/redis"
	"news_service/internal/repository"
	"news_service/internal/service"

	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	errorsUfb "github.com/leonardo849/utils_for_backend/pkg/errors"
	"github.com/thoas/go-funk"
)


var newsService *service.NewsService

func PermJournalist() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if newsService == nil {
			newsService = helper.CreateNewsService(repository.DB, redis.Rc)
		}
		mapClaims := ctx.Locals("user").(jwt.MapClaims)
		user := map[string]interface{}(mapClaims)
		authId := user["id"].(string)
		newsId := ctx.Params("news_id")
		
		status, message := newsService.FindAuthorsIdsByNews(newsId)
		if status >= 400 {
			return  ctx.Status(status).JSON(fiber.Map{"error": message})
		}
		authorsIds := message.([]string)
		if !funk.Contains(authorsIds, authId) {
			return ctx.Status(403).JSON(fmt.Sprintf("[%s] %s", errorsUfb.FORBIDDEN, "you can't acess that news"))
		}

		return ctx.Next()
	}
}