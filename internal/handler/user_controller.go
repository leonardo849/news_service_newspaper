package handler

import (
	"news_service/internal/dto"
	"news_service/internal/logger"
	"news_service/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type UserController struct {
	userService *service.UserService
}

func CreateNewUserController(userService *service.UserService) *UserController {
	return  &UserController{
		userService: userService,
	}
}

// @Summary     update bio
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       users body dto.UpdateAuthorBioDTO true "bio data"
// @Success     200 {object} dto.CreatedMessage
// @Failure     400 {object} dto.ErrorDTO
// @Failure     401 {object} dto.ErrorDTO
// @Failure     403 {object} dto.ErrorDTO
// @Failure     500 {object} dto.ErrorDTO
// @Router      /users/update/bio [patch]
// @Security    JWT
func (u *UserController) UpdateBio() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		mapClaims := ctx.Locals("user").(jwt.MapClaims)
		user := map[string]interface{}(mapClaims)
		authId := user["id"].(string)
		var input dto.UpdateAuthorBioDTO
		if err := ctx.BodyParser(&input); err != nil {
			logger.ZapLogger.Error("error in body parsing dto.UpdateAuthorBio", zap.Error(err))
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		status, message := u.userService.UpdateUserBio(authId, input)
		if status >= 400 {
			return  ctx.Status(400).JSON(fiber.Map{"error": message})
		}
		return  ctx.Status(200).JSON(message)
	}
}