package main

import (
	"log"
	"os"
	"template-backend/config"
	"template-backend/internal/logger"
	"template-backend/internal/prometheus"
	"template-backend/internal/repository"
	"template-backend/internal/router"
	"template-backend/internal/validate"

	_ "template-backend/docs"

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
	prometheus.StartPrometheus()
	validate.StartValidator()
	if err := router.RunServer(); err != nil {
		logger.ZapLogger.Error("error in run server", 
		zap.Error(err),
		zap.String("function", "router.RunServer()"),
		)
	}
}