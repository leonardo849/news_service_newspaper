package repository

import (
	"errors"
	"fmt"
	"news_service/internal/logger"
	"news_service/internal/model"

	"github.com/google/uuid"
	dtoSl "github.com/leonardo849/shared_library_news_paper/pkg/dto"
	errorsUfb "github.com/leonardo849/utils_for_backend/pkg/errors"
	"github.com/thoas/go-funk"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
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
	var searchedUser model.UserModel
	result := u.db.Where("username = ?", input.Username).First(&searchedUser)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.ZapLogger.Error("error finding user", zap.Error(result.Error))
		return nil, fmt.Errorf("[%s] error: %s", errorsUfb.INTERNALSERVER, result.Error.Error())
	}
	result = u.db.Create(&userModel)
	if result.Error != nil {
		logger.ZapLogger.Error("error creating user", zap.Error(result.Error))
		return nil, fmt.Errorf("[%s] error: %s", errorsUfb.INTERNALSERVER, result.Error.Error())
	}
	logger.ZapLogger.Info("user was created")
	return &userModel.ID, nil
}

func (u *UserRepository) CreateUsers(input []dtoSl.AuthPublishUserCreated) error {
	users := funk.Map(input, func(element dtoSl.AuthPublishUserCreated) *model.UserModel {
		return &model.UserModel{
			AuthId: element.AuthId,
			Role: element.Role,
			Username: element.Username,
		}
	}).([]*model.UserModel)
	
	result := u.db.Create(users)
	if result.Error != nil {
		logger.ZapLogger.Error("error in creating db users", zap.Error(result.Error))
		return  result.Error
	}
	logger.ZapLogger.Info("users were created")
	return  nil
}

func (u *UserRepository) SetDatabase(db *gorm.DB) {
	if u.db == nil {
		u.db = db
	}
}
