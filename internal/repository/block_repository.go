package repository

import (
	"errors"
	"fmt"
	"news_service/internal/logger"
	"news_service/internal/model"

	"github.com/google/uuid"
	errorsUfb "github.com/leonardo849/utils_for_backend/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BlockRepository struct {
	db *gorm.DB
}

func CreateBlockRepository(db *gorm.DB) *BlockRepository {
	return  &BlockRepository{
		db: db,
	}
}

func (b *BlockRepository) CreateBlock(tx *gorm.DB, input model.BlockModel) (*uuid.UUID, error) {
	if err := b.FindBlockByPosition(input.NewsID.String(), input.Position, tx); err != nil {
		logger.ZapLogger.Error("error find block by position", zap.Error(err))
		return nil, err
	}
	model := model.BlockModel{
		Content: input.Content,
		Position: input.Position,
		NewsID: input.NewsID,
	}

	

	if err := tx.Create(&model).Error; err != nil {
		logger.ZapLogger.Error("error creating block", zap.Error(err))
		return  nil, fmt.Errorf("[%s] %s", errorsUfb.INTERNALSERVER, err.Error())
	}
	return &model.ID, nil

}

func (b *BlockRepository) FindBlockByPosition(newsId string, position uint, tx *gorm.DB) error {
	var block model.BlockModel
	if err := tx.Where("news_id = ? AND position = ?", newsId, position).First(&block).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.ZapLogger.Error("error finding block", zap.Error(err))
			return  fmt.Errorf("[%s] %s", errorsUfb.INTERNALSERVER, err.Error())
		}
	}
	return nil
}