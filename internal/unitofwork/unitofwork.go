package unitofwork

import (
	"fmt"
	"news_service/internal/dto"
	"news_service/internal/logger"
	"news_service/internal/model"
	"news_service/internal/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UnitOfWork struct {
	userRepository *repository.UserRepository
	newsRepository *repository.NewsRepository
	imageRepository *repository.ImageRepository
	blockRepository *repository.BlockRepository
	db *gorm.DB
}

func CreateUnitOfWork(userRepository *repository.UserRepository, db *gorm.DB, imageRepository *repository.ImageRepository, newsRepository *repository.NewsRepository, blockRepository *repository.BlockRepository) *UnitOfWork {
	return  &UnitOfWork{
		userRepository: userRepository,
		newsRepository: newsRepository,
		imageRepository: imageRepository,
		blockRepository: blockRepository,
		db: db,
	}
}


func (u *UnitOfWork) createBlocksWithTx(input []dto.CreateBlockDTO, newsId uuid.UUID, tx *gorm.DB) error {
	for _, e := range input {
		blockModel := model.BlockModel{
			Content:  e.Content,
			Position: e.Position,
			NewsID:   newsId,
		}

		idBlock, err := u.blockRepository.CreateBlock(tx, blockModel)
		if err != nil {
			logger.ZapLogger.Error(fmt.Sprintf("error creating block. position %d", e.Position), zap.Error(err))
			return err
		}

		for _, i := range e.Images {
			imageModel := model.ImageModel{
				URL:     i.URL,
				BlockID: *idBlock,
			}

			if err := u.imageRepository.CreateImage(imageModel, tx); err != nil {
				logger.ZapLogger.Error("error creating image", zap.Error(err))
				return err
			}
		}
	}

	return nil
}


func (u *UnitOfWork) CreateNews(input dto.CreateNewsDTO, authors []model.UserModel) (string, error) {
	var id *uuid.UUID
	errTx := u.db.Transaction(func(tx *gorm.DB) error {
		var err error
		id, err = u.newsRepository.CreateNews(input, authors, tx)
		if err != nil {
			return err
		}

		if err := u.createBlocksWithTx(input.Blocks, *id, tx); err != nil {
			return err
		}

		return nil
	})

	if id != nil {
		return id.String(), errTx
	}
	return "", errTx
}


func (u *UnitOfWork) CreateBlocks(input []dto.CreateBlockDTO, newsId uuid.UUID) error {
	errTx := u.db.Transaction(func(tx *gorm.DB) error {
		return u.createBlocksWithTx(input, newsId, tx)
	})

	if errTx == nil {
		logger.ZapLogger.Info("blocks were created")
	}
	return errTx
}
