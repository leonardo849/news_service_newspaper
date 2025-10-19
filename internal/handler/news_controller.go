package handler

import (
	"news_service/internal/dto"
	_ "news_service/internal/dto"
	"news_service/internal/logger"
	_ "news_service/internal/repository"
	"news_service/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type NewsController struct {
	newsService *service.NewsService
}

func CreateNewsController(newsService *service.NewsService) *NewsController {
	return &NewsController{
		newsService: newsService,
	}
}

// @Summary     create news
// @Tags        news
// @Accept      json
// @Produce     json
// @Param       news body dto.CreateNewsDTO true "news data"
// @Success     200 {object} dto.CreatedMessage
// @Failure     400 {object} dto.ErrorDTO
// @Failure     401 {object} dto.ErrorDTO
// @Failure     403 {object} dto.ErrorDTO
// @Failure     500 {object} dto.ErrorDTO
// @Router      /news/create [post]
// @Security    JWT
func (n *NewsController) CreateNews() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		mapClaims := ctx.Locals("user").(jwt.MapClaims)
		user := map[string]interface{}(mapClaims)
		authId := user["id"].(string)
		var input dto.CreateNewsDTO
		if err := ctx.BodyParser(&input); err != nil {
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		status, message := n.newsService.CreateNews(input, authId)
		if status >= 400 {
			return ctx.Status(status).JSON(fiber.Map{"error": message})
		}
		return ctx.Status(200).JSON(message)
	}
}

// @Summary     find news
// @Tags        news
// @Accept      json
// @Produce     json
// @Param        id   path      string  true  "news ID"
// @Success     200 {object} dto.FindNewsDTO
// @Failure     400 {object} dto.ErrorDTO
// @Failure     401 {object} dto.ErrorDTO
// @Failure     403 {object} dto.ErrorDTO
// @Failure     500 {object} dto.ErrorDTO
// @Router      /news/{id} [get]
// @Security    JWT
func (n *NewsController) FindNewsById() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		status, message := n.newsService.FindNewsById(id)
		if status >= 400 {
			logger.ZapLogger.Error("error " + message.(string))
			return  ctx.Status(status).JSON(fiber.Map{"error": message})
		}
		return  ctx.Status(200).JSON(message)
	}
}

// @Summary     publish news
// @Tags        news
// @Accept      json
// @Produce     json
// @Param        id   path      string  true  "news ID"
// @Success     200 {object} dto.MessageDTO
// @Failure     400 {object} dto.ErrorDTO
// @Failure     401 {object} dto.ErrorDTO
// @Failure     403 {object} dto.ErrorDTO
// @Failure     500 {object} dto.ErrorDTO
// @Router      /news/publish/{id} [patch]
// @Security    JWT
func (n *NewsController) PublishNewsById() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		status, message := n.newsService.PublishNews(id)
		if status >= 400 {
			logger.ZapLogger.Error("error " + message.(string))
			return  ctx.Status(status).JSON(fiber.Map{"error": message})
		}
		return ctx.Status(200).JSON(message)
	}
}