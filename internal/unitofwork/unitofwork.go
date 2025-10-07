package unitofwork

import "news_service/internal/repository"

type UnitOfWork struct {
	userRepository *repository.UserRepository
}

func CreateUnitOfWork(userRepository *repository.UserRepository) *UnitOfWork {
	return  &UnitOfWork{
		userRepository: userRepository,
	}
}

