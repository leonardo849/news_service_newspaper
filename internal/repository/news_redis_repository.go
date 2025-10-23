package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"news_service/internal/dto"
	"news_service/internal/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)


type NewsRedisRepository struct {
	rc    *redis.Client
	model string
}


func CreateNewsRedisRepository(rc *redis.Client) *NewsRedisRepository {
	return &NewsRedisRepository{
		rc:    rc,
		model: "news",
	}
}


func (n *NewsRedisRepository) SetNews(input dto.FindNewsDTO, ctx context.Context) error {
	key := n.getRedisKey(input.ID)

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		logger.ZapLogger.Error("error dto.FindNewsDTO to bytes", zap.Error(err))
		return err
	}

	redisStatus := n.rc.Set(ctx, key, jsonBytes, 10*time.Minute)
	if redisStatus.Err() != nil {
		logger.ZapLogger.Error("error setting news in Redis", zap.Error(redisStatus.Err()))
		return redisStatus.Err()
	}

	logger.ZapLogger.Info(fmt.Sprintf("news was setted in redis key: %s", key))
	return nil
}


func (n *NewsRedisRepository) GetNews(id string, ctx context.Context) (*dto.FindNewsDTO, error) {
	key := n.getRedisKey(id)

	val, err := n.rc.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			logger.ZapLogger.Error(fmt.Sprintf("news wasn't found in Redis. key: %s", key))
			return nil, err
		}
		logger.ZapLogger.Error("error getting news from Redis", zap.Error(err))
		return nil, err
	}

	var news dto.FindNewsDTO
	if err := json.Unmarshal([]byte(val), &news); err != nil {
		logger.ZapLogger.Error("error bytes to json", zap.Error(err))
		return nil, err
	}

	logger.ZapLogger.Info(fmt.Sprintf("news was gotten from Redis. key: %s", key))
	return &news, nil
}


func (n *NewsRedisRepository) getRedisKey(id string) string {
	return fmt.Sprintf("%s:%s", n.model, id)
}
