package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"news_service/internal/dto"
	"news_service/internal/logger"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type NewsRedisRepository struct {
	rc *redis.Client
	model string
}


func getRedisKey(id string, model string) string {
	return  model+":"+id
}

func CreateNewsRedisRepository(rc *redis.Client) *NewsRedisRepository {
	return  &NewsRedisRepository{
		rc: rc,
		model: "news",
	}
}

func (n *NewsRedisRepository) SetNews(input dto.FindNewsDTO, context context.Context) error {
	key := getRedisKey(input.ID, n.model)
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		logger.ZapLogger.Error("error in json marshal input dto.Findnewsdto", zap.Error(err))
		return err
	}
	if redisStatus := n.rc.Set(context, key, jsonBytes, 10 * time.Minute); redisStatus.Err() != nil {
		logger.ZapLogger.Error("error in set input dto.Findnewsdto in database", zap.Error(err))
		return  err
	}
	logger.ZapLogger.Info(fmt.Sprintf("news was setted in redis. id: %s", input.ID))
	
	return  nil
} 

func (n *NewsRedisRepository) GetNews(id string, context context.Context) (*dto.FindNewsDTO, error) {
	key := getRedisKey(id, n.model)
	val, err := n.rc.Get(context, key).Result()
	if err != nil {
		logger.ZapLogger.Error("error in get news", zap.Error(err))
		return nil, err
	}
	var news dto.FindNewsDTO
	err = json.Unmarshal([]byte(val), &news)
	if err != nil {
		logger.ZapLogger.Error("error in json unmarshal", zap.Error(err))
		return nil, err
	}

	logger.ZapLogger.Info(fmt.Sprintf("news with id %s was gotten from redis", id))
	return  &news, nil
}