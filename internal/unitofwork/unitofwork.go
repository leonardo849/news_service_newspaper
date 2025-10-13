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

func (u *UnitOfWork) CreateNews(input dto.CreateNewsDTO, authors []model.UserModel) (string, error) {
	var id *uuid.UUID
	err := u.db.Transaction(func(tx *gorm.DB) error {
		var err error
		
		
		

		id, err = u.newsRepository.CreateNews(input, authors, tx)
		if err != nil {
			return err
		}
		
		for _, e := range input.Blocks {
			blockModel := model.BlockModel{
				Content: e.Content,
				Position: e.Position,
				NewsID: *id,
			}
			var idBlock *uuid.UUID
			if idBlock, err = u.blockRepository.CreateBlock(tx, blockModel); err != nil {
				logger.ZapLogger.Error(fmt.Sprintf("error creating block. position %d", e.Position), zap.Error(err))
				return err
			}
			for _, i := range e.Images {
				imageModel := model.ImageModel{
					URL: i.URL,
					BlockID: *idBlock,
				}
				if err = u.imageRepository.CreateImage(imageModel, tx); err != nil {
					logger.ZapLogger.Error("error creating image", zap.Error(err))
					return  err
				}
			}
		}

		

		
		return  nil
	})
	return  id.String(), err
}

