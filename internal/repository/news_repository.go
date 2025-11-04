package repository

import (
	"errors"
	"fmt"
	"news_service/internal/dto"
	"news_service/internal/helper_consts"
	"news_service/internal/logger"
	"news_service/internal/model"
	"time"

	"github.com/google/uuid"
	"github.com/leonardo849/utils_for_backend/pkg/date"
	errorsUfb "github.com/leonardo849/utils_for_backend/pkg/errors"
	"github.com/thoas/go-funk"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NewsRepository struct {
	db *gorm.DB
}

func CreateNewsRepository(db *gorm.DB) *NewsRepository {
	return  &NewsRepository{
		db: db,
	}
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

func (n *NewsRepository) PublishNews(id uuid.UUID) error {
	if err := n.db.Model(&model.NewsModel{}).Where("id = ? AND status = ?", id, helper_consts.SKETCH).Updates(map[string]interface{}{"status": helper_consts.PUBLISHED, "published_at": date.PtrTime(time.Now())}).Error; err != nil {
		return fmt.Errorf("[%s] %s", errorsUfb.INTERNALSERVER, err.Error())
	}
	return  nil
}

func (n *NewsRepository) FindNewsById(id uuid.UUID) (*model.NewsModel, error) {
	var news model.NewsModel
	if err := n.db.Where("id = ? AND status = ?", id, helper_consts.PUBLISHED).Preload("Authors").Preload("Blocks", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC")
	}).
	Preload("Blocks.Images").First(&news).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.ZapLogger.Error("error", zap.Error(err))
			return  nil, fmt.Errorf("[%s] %s", errorsUfb.NOTFOUND, "news wasn't found")
		} else {
			logger.ZapLogger.Error("error", zap.Error(err))
			return  nil, fmt.Errorf("[%s] %s", errorsUfb.INTERNALSERVER, err.Error())
		}
	}
	return  &news, nil

}

func (n *NewsRepository) FindNotPublishedNews(id uuid.UUID) (*model.NewsModel, error) {
	var news model.NewsModel
	if err := n.db.Where("id = ? AND status = ?", id, helper_consts.SKETCH).Preload("Authors").Preload("Blocks", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC")}).First(&news).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.ZapLogger.Error("error", zap.Error(err))
			return  nil, fmt.Errorf("[%s] %s", errorsUfb.NOTFOUND, "news wasn't found")
		} else {
			logger.ZapLogger.Error("error", zap.Error(err))
			return  nil, fmt.Errorf("[%s] %s", errorsUfb.INTERNALSERVER, err.Error())
		}
	}
	return  &news, nil
}

func (n *NewsRepository) FindAuthorsIdsByNewsId(newsId uuid.UUID) ([]string, error) {
	news, err := n.FindNotPublishedNews(newsId)
	if err != nil {
		logger.ZapLogger.Error("error", zap.Error(err))
		return nil, err
	}
	ids := funk.Map(news.Authors, func(author model.UserModel) string {
		return author.AuthId
	}).([]string)

	return ids, nil
}