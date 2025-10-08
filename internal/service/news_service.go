package service

import (
	"news_service/internal/dto"
	"news_service/internal/repository"
	"news_service/internal/validate"
	errorsSl "github.com/leonardo849/shared_library_news_paper/pkg/errors"
	"github.com/google/uuid"
)

type NewsService struct {
	newsRepository *repository.NewsRepository
	userRepository *repository.UserRepository
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
		status, message = errorsSl.HandleErrors(err, n.model)
		return status, message
	}

	n.newsRepository.CreateNews(input, authors)


}