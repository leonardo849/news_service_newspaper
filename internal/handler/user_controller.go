package handler

import (
	"news_service/internal/service"

)

type UserController struct {
	userService *service.UserService
}

func CreateNewUserController(userService *service.UserService) *UserController {
	return  &UserController{
		userService: userService,
	}
}

