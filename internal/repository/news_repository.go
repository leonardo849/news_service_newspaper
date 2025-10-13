package repository

import (
	"news_service/internal/dto"
	"news_service/internal/logger"
	"news_service/internal/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NewsRepository struct {
	db *gorm.DB
}

func (n *NewsRepository) CreateNews(input dto.CreateNewsDTO, authors []model.UserModel, tx *gorm.DB) (*uuid.UUID, error) {
	news := model.NewsModel{
		Title:input.Title,
		Subtitle: input.Subtitle,
		Topic: input.Topic,
		Authors: authors,
	}
	err := tx.Create(&news).Error
	if err != nil {
		logger.ZapLogger.Error("error creating news", zap.Error(err))
		return nil, err
	}
	return &news.ID, nil
}