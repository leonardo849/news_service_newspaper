package repository

import (
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
	if b.IsThereABlockInPosition(int(input.Position), input.NewsID.String(), tx) {
		return nil, fmt.Errorf("[%s] %s", errorsUfb.CONFLICT, "a block in this position")
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

func (b *BlockRepository) IsThereABlockInPosition(position int, newsId string, tx *gorm.DB) bool {
	var block model.BlockModel
	err := tx.Where("news_id = ? AND position = ?", newsId, position).First(&block).Error
	return err == nil
}