package handler

import (
	"news_service/internal/dto"
	"news_service/internal/logger"
	"news_service/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	// "github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type BlockController struct {
	blockService *service.BlockService
}

func CreateBlockController(blockService *service.BlockService) *BlockController {
	return &BlockController{
		blockService: blockService,
	}
}

// @Summary     create blocks
// @Tags        news
// @Accept      json
// @Produce     json
// @Param        news_id  path      string  true  "news ID"
// @Param       news body dto.CreateBlockDTO true "block data"
// @Success     200 {object} dto.MessageDTO
// @Failure     400 {object} dto.ErrorDTO
// @Failure     401 {object} dto.ErrorDTO
// @Failure     403 {object} dto.ErrorDTO
// @Failure     500 {object} dto.ErrorDTO
// @Router      /blocks/create/{news_id} [post]
// @Security    JWT
func (b *BlockController) CreateBlocks() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		var input []dto.CreateBlockDTO
		if err := ctx.BodyParser(&input); err != nil {
			logger.ZapLogger.Error("error in body parser dto.createblockdto", zap.Error(err))
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		mapClaims := ctx.Locals("user").(jwt.MapClaims)
		user := map[string]interface{}(mapClaims)
		authId := user["id"].(string)
		newsId := ctx.Params("news_id")

		status, messsage := b.blockService.CreateBlocks(input, newsId, authId)
		if status >= 400 {
			return ctx.Status(status).JSON(fiber.Map{"error": messsage})
		}
		return ctx.Status(200).JSON(messsage)
	}
}
