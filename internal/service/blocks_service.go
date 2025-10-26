package service

import (
	"fmt"
	"news_service/internal/dto"
	"news_service/internal/logger"
	"news_service/internal/unitofwork"

	"github.com/google/uuid"
	errorsSl "github.com/leonardo849/shared_library_news_paper/pkg/errors"
	errorsUfb "github.com/leonardo849/utils_for_backend/pkg/errors"
	"go.uber.org/zap"
)

type BlockService struct {
	newsService      *NewsService
	unitOfWork          *unitofwork.UnitOfWork
	model               string
}

func CreateBlockService(newsService *NewsService,  unitOfWork *unitofwork.UnitOfWork) *BlockService {
	return  &BlockService{
		newsService: newsService,
		unitOfWork: unitOfWork,
		model: "blolcks",
	}
}



func (b *BlockService) CreateBlocks(input []dto.CreateBlockDTO, newsId string, authId string) (status int, message interface{}) {
	newsIdToUuid, err := uuid.Parse(newsId)
	if err != nil {
		logger.ZapLogger.Error("error", zap.Error(err))
		return 500, fmt.Errorf("[%s] %s", errorsUfb.INTERNALSERVER, err.Error())
	}
	status, message = b.newsService.FindNotPublishedNews(newsId, authId)
	if status >= 400 {
		return status, message
	}
	if err := b.unitOfWork.CreateBlocks(input, newsIdToUuid); err != nil {
		logger.ZapLogger.Error("error creating blocks", zap.Error(err))
		status, message = errorsSl.HandleErrors(err, b.model)
		return status, message
	}


	

	return status, dto.MessageDTO{
		Message: "blocks were created",
	}
}