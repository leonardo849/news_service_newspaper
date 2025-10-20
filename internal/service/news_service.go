package service

import (
	"news_service/internal/dto"
	"news_service/internal/logger"
	"news_service/internal/model"
	"news_service/internal/repository"
	"news_service/internal/unitofwork"
	"news_service/internal/validate"
	

	"github.com/google/uuid"
	errorsSl "github.com/leonardo849/shared_library_news_paper/pkg/errors"
	"github.com/thoas/go-funk"
	"go.uber.org/zap"
)

type NewsService struct {
	newsRepository *repository.NewsRepository
	userRepository *repository.UserRepository
	newsRedisRepository *repository.NewsRedisRepository
	unitOfWork *unitofwork.UnitOfWork
	model string
}

func CreateNewsService(newsRepository *repository.NewsRepository, userRepository *repository.UserRepository, unitOfWork *unitofwork.UnitOfWork, newsRedisRepository *repository.NewsRedisRepository ) *NewsService {
	return  &NewsService{
		newsRepository: newsRepository,
		userRepository: userRepository,
		unitOfWork: unitOfWork,
		model: "news",
		newsRedisRepository: newsRedisRepository,
	}
}

func (n *NewsService) CreateNews(input dto.CreateNewsDTO, id string) (status int, message interface{}) {
	if err := validate.Validate.Struct(input); err != nil {
		return 400, err.Error()
	}
	
	idUser, err := n.userRepository.FindOneUserIdByAuthId(id)
	if err != nil {
		logger.ZapLogger.Error("error userrepository.findoneuseridbyauthid", zap.Error(err))
		status, message = errorsSl.HandleErrors(err, n.model)
		return status, message
	}

	ids := []uuid.UUID{*idUser}

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

func (n *NewsService) FindNewsById(id string) (status int, message interface{}) {
	news, err := n.newsRepository.FindNewsById(id)
	if err != nil {
		logger.ZapLogger.Error("error finding news by id " + news.ID.String(), zap.Error(err))
		status, message = errorsSl.HandleErrors(err, n.model)
		return status, message
	}
	authors := funk.Map(news.Authors, func(a model.UserModel) dto.FindAuthorsInFindNews {
		return dto.FindAuthorsInFindNews{
			ID: a.ID.String(),
			Username: a.Username,
		}
	}).([]dto.FindAuthorsInFindNews)
	blocks := funk.Map(news.Blocks, func(b model.BlockModel) dto.FindBlockDTO {
	images := funk.Map(b.Images, func(i model.ImageModel) dto.FindImageDTO {
		return dto.FindImageDTO{
			ID: i.ID.String(),
			URL:     i.URL,
		}
	}).([]dto.FindImageDTO) 

		return dto.FindBlockDTO{
			ID: b.ID.String(),
			Position: b.Position,
			Content:  b.Content,
			Images:   images,
		}
	}).([]dto.FindBlockDTO) 

	

	newsDto := dto.FindNewsDTO{
		ID: news.ID.String(),
		Authors: authors,
		Title: news.Title,
		Subtitle: news.Subtitle,
		Topic: news.Topic,
		Blocks: blocks,
		CreatedAt: news.CreatedAt,
		UpdatedAt: news.UpdatedAt,
		Published_at: news.Published_at,
	}
	return 200, newsDto
}

func (n *NewsService) PublishNews(id string) (status int, message interface{}) {
	if err := n.newsRepository.PublishNews(id); err != nil {
		logger.ZapLogger.Error("error publishing news by id", zap.Error(err))
		status, message = errorsSl.HandleErrors(err, n.model)
		return status, message
	}
	return 200, dto.MessageDTO{Message: "news was published"}
}