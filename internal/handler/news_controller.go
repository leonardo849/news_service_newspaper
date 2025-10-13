package handler

import (
	_"news_service/internal/dto"
	_"news_service/internal/repository"
	"news_service/internal/service"

	_"github.com/gofiber/fiber/v2"
)

type NewsController struct {
	newsService *service.NewsService
}

func CreateNewNewsController(newsService *service.NewsService) *NewsController {
	return &NewsController{
		newsService: newsService,
	}
}

// func (n *NewsController) CreateNews() fiber.Handler {
// 	return func(ctx *fiber.Ctx) error {
// 		var input dto.CreateNewsDTO
// 		if err := ctx.BodyParser(&input); err != nil {
// 			return  ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
// 		}

// 		n.newsService.CreateNews(input,)
// 	}
// }