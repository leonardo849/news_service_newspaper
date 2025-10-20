package service

import (
	"news_service/internal/dto"
	"news_service/internal/logger"
	"news_service/internal/repository"
	"news_service/internal/validate"

	dtoSl "github.com/leonardo849/shared_library_news_paper/pkg/dto"
	errorsSl "github.com/leonardo849/shared_library_news_paper/pkg/errors"
	"go.uber.org/zap"
)

type UserService struct {
	userRepository *repository.UserRepository
	model          string
}

func CreateNewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
		model:          "user",
	}
}

func (u *UserService) CreateUser(input dtoSl.AuthPublishUserCreated) (status int, message interface{}) {
	id, err := u.userRepository.CreateUser(input)
	if err != nil {
		status, message = errorsSl.HandleErrors(err, u.model)
		return status, message
	}

	idStr := id.String()

	return 200, map[string]string{
		"id":      idStr,
		"message": "user was generated",
	}
}

func (u *UserService) FindUserByUsername(username string) (status int, message interface{}) {
	user, err := u.userRepository.FindUserByUsername(username)
	if err != nil {
		status, message = errorsSl.HandleErrors(err, u.model)
		return status, message
	}
	return 200, *user
}

func (u *UserService) CreateUsers(input []dtoSl.AuthPublishUserCreated) (status int, message string) {
	err := u.userRepository.CreateUsers(input)
	if err != nil {
		status, message = errorsSl.HandleErrors(err, u.model)
		return status, message
	}
	return 200, "users were created"
}


func (u *UserService)UpdateUserBio(authId string, input dto.UpdateAuthorBioDTO) (status int, message interface{}) {
	if err := validate.Validate.Struct(input); err != nil {
		logger.ZapLogger.Error("error validating struct dto.UpdateAuthorBio", zap.Error(err))
		return 400, err.Error()
	}
	if _, err := u.userRepository.FindOneUserIdByAuthId(authId); err != nil {
		logger.ZapLogger.Error("error in u.userRepository.FindOneUserIdByAuthId", zap.Error(err))
		status, message = errorsSl.HandleErrors(err, u.model)
		return status, message
	}
	if err := u.userRepository.UpdateBio(input.Bio, authId); err != nil {
		logger.ZapLogger.Error("error in u.userRepository.UpdateBio", zap.Error(err))
		status, message = errorsSl.HandleErrors(err, u.model)
		return 500, err.Error()
	}
	return 200, dto.MessageDTO{
		Message: "user's bio was updated",
	}
}