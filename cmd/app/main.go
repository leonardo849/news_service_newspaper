package main

import (
	"log"
	"os"
	"news_service/config"
	"news_service/internal/logger"
	"news_service/internal/prometheus"
	"news_service/internal/redis"
	"news_service/internal/repository"
	"news_service/internal/router"
	"news_service/internal/validate"

	_ "news_service/docs"

	"go.uber.org/zap"
)

// @title Backend Portfolio API
// @version 1.0
// @description api for a news service
// @host localhost:port
// @BasePath /
func main() {
	if err := config.SetupEnvVar(); err != nil {
		log.Fatal(err.Error())
	}
	if err := logger.StartLogger(); err != nil {
		log.Fatal(err.Error())
	}
	if _,err := repository.ConnectToDatabase(); err != nil {
		logger.ZapLogger.Error("error in repository.connectodatabase", zap.String("function", "repository.ConnectToDatabase()"), zap.Error(err))
		os.Exit(1)
	}
	if _, err := redis.ConnectToRedis(); err != nil {
		logger.ZapLogger.Error("error in connect to redis", zap.String("function", "redis.ConnectToRedis"), zap.Error(err))
		os.Exit(1)
	}
	prometheus.StartPrometheus()
	validate.StartValidator()
	if err := router.RunServer(); err != nil {
		logger.ZapLogger.Error("error in run server", 
		zap.Error(err),
		zap.String("function", "router.RunServer()"),
		)
	}
}