package repository

import (
	"errors"
	"fmt"
	"news_service/internal/logger"
	"news_service/internal/model"

	"github.com/google/uuid"
	constsSl "github.com/leonardo849/shared_library_news_paper/pkg/consts"
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
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) CreateUser(input dtoSl.AuthPublishUserCreated) (*uuid.UUID, error) {
	userModel := model.UserModel{
		Username: input.Username,
		Role:     input.Role,
		AuthId:   input.AuthId,
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

func (u *UserRepository) FindOneUserIdByAuthId(authId string) (*uuid.UUID, error) {
	var user model.UserModel
	if err := u.db.Where("auth_id = ?", authId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("[%s] %s", errorsUfb.NOTFOUND, "user wasn't found")
		} else {
			return nil, fmt.Errorf("[%s] %s", errorsUfb.INTERNALSERVER, err.Error())
		}
	}
	return &user.ID, nil
}



func (u *UserRepository) CreateUsers(input []dtoSl.AuthPublishUserCreated) error {
	

	users := funk.Map(input, func(element dtoSl.AuthPublishUserCreated) *model.UserModel {
		return &model.UserModel{
			AuthId:   element.AuthId,
			Role:     element.Role,
			Username: element.Username,
		}
	}).([]*model.UserModel)



	result := u.db.Create(users)
	if result.Error != nil {
		logger.ZapLogger.Error("error in creating db users", zap.Error(result.Error))
		return result.Error
	}
	logger.ZapLogger.Info("users were created")
	return nil
}

func (u *UserRepository) DoAuthorsExistAndReturnAuthors(ids []uuid.UUID) ([]model.UserModel, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("[%s]there is no id", errorsUfb.BADREQUEST)
	}
	var count int64
	var users []model.UserModel
	err := u.db.Where("id IN ?", ids).Find(&users).Count(&count).Error
	if err != nil {
		logger.ZapLogger.Error("error in find authors", zap.Error(err))
		return nil, fmt.Errorf("[%s] %s", errorsUfb.NOTFOUND, err.Error())
	}
	for _, user := range users {
		if user.Role != constsSl.Journalist {
			return nil, fmt.Errorf("[%s]users aren't journalist ", errorsUfb.UNAUTHORIZED)
		}
	}
	return users, nil
}

func (u *UserRepository) FindUserByUsername(username string) (*model.UserModel, error) {
	var user model.UserModel
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		logger.ZapLogger.Error("error", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("[%s] %s", errorsUfb.NOTFOUND, "user wasn't found by username "+username)
		} else {
			return nil, fmt.Errorf("[%s] %s", errorsUfb.INTERNALSERVER, err.Error())
		}
	}
	return &user, nil
}

func (u *UserRepository) SetDatabase(db *gorm.DB) {
	if u.db == nil {
		u.db = db
	}
}
