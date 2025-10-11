package repository

import (
	"errors"
	"fmt"
	"news_service/internal/logger"
	"news_service/internal/model"

	errorsUfb "github.com/leonardo849/utils_for_backend/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BlockRepository struct {
	db *gorm.DB
}

func (b *BlockRepository) CreateBlock(tx *gorm.DB, input model.BlockModel) error {
	if err := b.FindBlockByPosition(input.NewsID.String(), input.Position); err != nil {
		logger.ZapLogger.Error("error find block by position", zap.Error(err))
		return err
	}
	model := &model.BlockModel{
		Content: input.Content,
		Position: input.Position,
		NewsID: input.NewsID,
	}

	

	if err := b.db.Create(&model).Error; err != nil {
		logger.ZapLogger.Error("error creating block", zap.Error(err))
		return  fmt.Errorf("[%s] %s", errorsUfb.INTERNALSERVER, err.Error())
	}
	return nil

}

func (b *BlockRepository) FindBlockByPosition(newsId string, position uint) error {
	var block model.BlockModel
	if err := b.db.Where("news_id = ? AND position = ?", newsId, position).First(&block).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.ZapLogger.Error("error finding block", zap.Error(err))
			return  fmt.Errorf("[%s] %s", errorsUfb.INTERNALSERVER, err.Error())
		}
	}
	return nil
}