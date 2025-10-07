package main

import (
	"log"
	"news_service/config"
	"news_service/internal/logger"
	"news_service/internal/prometheus"
	"news_service/internal/rabbitmq"
	"news_service/internal/redis"
	"news_service/internal/repository"
	"news_service/internal/router"
	"news_service/internal/validate"
	"os"

	_ "news_service/docs"

	"go.uber.org/zap"
)

// @title Backend News service
// @version 1.0
// @description api for a news service
// @host localhost:port
// @BasePath /
func main() {
	if err := config.SetupEnvVar(); err != nil {
		log.Fatal(err.Error())
	}
	log.Print("env vars are ready")
	if err := logger.StartLogger(); err != nil {
		log.Fatal(err.Error())
	}
	logger.ZapLogger.Info("logger is ready")
	if err := rabbitmq.ConnectToRabbitMQ(); err != nil {
		logger.ZapLogger.Error(err.Error(), zap.Error(err))
	}
	logger.ZapLogger.Info("rabbit is ready")
	if _,err := repository.ConnectToDatabase(); err != nil {
		logger.ZapLogger.Error("error in repository.connectodatabase", zap.String("function", "repository.ConnectToDatabase()"), zap.Error(err))
		os.Exit(1)
	}
	logger.ZapLogger.Info("db is ready")
	if _, err := redis.ConnectToRedis(); err != nil {
		logger.ZapLogger.Error("error in connect to redis", zap.String("function", "redis.ConnectToRedis"), zap.Error(err))
		os.Exit(1)
	}
	logger.ZapLogger.Info("redis is ready")
	prometheus.StartPrometheus()
	logger.ZapLogger.Info("prometheus is ready")
	validate.StartValidator()
	logger.ZapLogger.Info("validate is ready")
	if err := router.RunServer(); err != nil {
		logger.ZapLogger.Error("error in run server", 
		zap.Error(err),
		zap.String("function", "router.RunServer()"),
		)
	}
}