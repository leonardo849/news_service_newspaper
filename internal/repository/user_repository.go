package repository

import (
	"news_service/internal/logger"
	"news_service/internal/model"

	"github.com/google/uuid"
	dtoSl "github.com/leonardo849/shared_library_news_paper/pkg/dto"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db*gorm.DB
}

func CreateUserRepository(db *gorm.DB) *UserRepository {
	return  &UserRepository{
		db: db,
	}
}

func (u *UserRepository) CreateUser(input dtoSl.AuthPublishUserCreated) (*uuid.UUID, error){
	userModel := model.UserModel{
		Username: input.Username,
		Role: input.Role,
		AuthId: input.AuthId,
	}
	result := u.db.Create(&userModel)
	if result.Error != nil {
		logger.ZapLogger.Error("error creating user", zap.Error(result.Error))
		return nil, result.Error
	}
	logger.ZapLogger.Info("user was created")
	return &userModel.ID, nil
}