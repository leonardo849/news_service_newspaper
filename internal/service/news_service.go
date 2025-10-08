package service

import (
	"news_service/internal/dto"
	"news_service/internal/logger"
	"news_service/internal/repository"
	"news_service/internal/unitofwork"
	"news_service/internal/validate"

	"github.com/google/uuid"
	errorsSl "github.com/leonardo849/shared_library_news_paper/pkg/errors"
	"go.uber.org/zap"
)

type NewsService struct {
	newsRepository *repository.NewsRepository
	userRepository *repository.UserRepository
	unitOfWork *unitofwork.UnitOfWork
	model string
}

func (n *NewsService) CreateNews(input dto.CreateNewsDTO, id string) (status int, message interface{}) {
	if err := validate.Validate.Struct(input); err != nil {
		return 400, err.Error()
	}
	ids := []uuid.UUID{uuid.MustParse(id)}

	for _, e := range input.Authors {
		ids = append(ids, uuid.MustParse(e))
	}

	authors, err := n.userRepository.DoAuthorsExistAndReturnAuthors(ids)
	if err != nil {
		logger.ZapLogger.Error("error userRepository.DoAuthorsExistAndReturnAuthors", zap.Error(err))
		status, message = errorsSl.HandleErrors(err, n.model)
		return status, message
	}

	newsId, err := n.unitOfWork.CreateNews(input, authors)
	if err != nil {
		logger.ZapLogger.Error("error unitOfWork.CreateNews", zap.Error(err))
		status, message = errorsSl.HandleErrors(err, n.model)
		return status, message
	}

	return 200, map[string]string{
		"id": newsId,
		"message": "a new news was generated",
	}
}