package unitofwork

import (
	"news_service/internal/dto"
	"news_service/internal/model"
	"news_service/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UnitOfWork struct {
	userRepository *repository.UserRepository
	newsRepository *repository.NewsRepository
	db *gorm.DB
}

func CreateUnitOfWork(userRepository *repository.UserRepository, db *gorm.DB) *UnitOfWork {
	return  &UnitOfWork{
		userRepository: userRepository,
		db: db,
	}
}

func (u *UnitOfWork) CreateNews(input dto.CreateNewsDTO, authors []model.UserModel) (string, error) {
	var id *uuid.UUID
	err := u.db.Transaction(func(tx *gorm.DB) error {
		var err error
		

		id, err = u.newsRepository.CreateNews(input, authors, tx)
		if err != nil {
			return nil
		}
		return  nil
	})
	return  id.String(), err
}

