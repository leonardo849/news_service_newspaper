package unitofwork

import "news_service/internal/repository"

type UnitOfWork struct {
	userRepository *repository.UserRepository
}

