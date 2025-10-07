package service

import (
	"news_service/internal/repository"

	dtoSl "github.com/leonardo849/shared_library_news_paper/pkg/dto"
	errorsSl "github.com/leonardo849/shared_library_news_paper/pkg/errors"
)

type UserService struct {
	userRepository *repository.UserRepository
	model string
}

func CreateNewUserService(userRepository *repository.UserRepository) *UserService {
	return  &UserService{
		userRepository: userRepository,
		model: "user",
	}
}

func (u *UserService) CreateUser(input dtoSl.AuthPublishUserCreated) (status int, message interface{}){
	id, err := u.userRepository.CreateUser(input)
	if err != nil {
		status, message = errorsSl.HandleErrors(err, u.model)
		return status, message
	}
	
	idStr := id.String()

	return 200, map[string]string{
		"id": idStr,
		"message": "user was generated",
	} 
}

func (u *UserService) CreateUsers(input []dtoSl.AuthPublishUserCreated) (status int, message string) {
	err := u.userRepository.CreateUsers(input)
	if err != nil {
		return 500, err.Error()
	}
	return 200, "users were created"
}