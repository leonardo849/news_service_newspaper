package repository

import (
	"fmt"
	"news_service/internal/helper"
	"news_service/internal/logger"
	"news_service/internal/model"

	errorsUfb "github.com/leonardo849/utils_for_backend/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ImageRepository struct {
	db *gorm.DB
}

func (i *ImageRepository) CreateImage(input model.ImageModel, tx *gorm.DB) error {
	helper.SetTx(&tx, i.db)
	image := model.ImageModel{
		URL: input.URL,
		BlockID: input.BlockID,
	}
	if err := i.db.Create(&image).Error; err != nil {
		logger.ZapLogger.Error("error creating image", zap.Error(err))
		return fmt.Errorf("[%s] %s", errorsUfb.INTERNALSERVER, err.Error())
	}
	return  nil
}