package integration_test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"news_service/config"
	"news_service/internal/dto"
	"news_service/internal/logger"
	"news_service/internal/model"
	_ "news_service/internal/model"
	"news_service/internal/rabbitmq"
	"news_service/internal/redis"
	"news_service/internal/repository"
	"news_service/internal/router"
	"news_service/internal/validate"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	redisLib "github.com/redis/go-redis/v9"

	"github.com/gavv/httpexpect/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var app *fiber.App


type fiberRoundTripper struct {
	app *fiber.App
}





func (rt fiberRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.app.Test(req)
}

var DB *gorm.DB

func TestMain(m *testing.M) {
	err := config.SetupEnvVar()
	if err != nil {
		log.Panic(err.Error())
	}
	os.Setenv("RABBIT_ON", "false")
	if err = logger.StartLogger(); err != nil {
		log.Panic(err.Error())
	}
	db, err := repository.ConnectToDatabase()
	if err != nil {
		log.Panic(err.Error())
	}
	rc, err := redis.ConnectToRedis()
	if err != nil {
		log.Panic(err.Error())
	}

	err = rabbitmq.ConnectToRabbitMQ()
	if err != nil {
		log.Print(err.Error())
	}
	validate.StartValidator()
	
	DB = db
	app = router.SetupApp(db, rc)
	sqldb, err := db.DB()
	if err != nil {
		log.Panic(err.Error())
	}
	
	cleanDatabase(db, rc)
	migrateSeeds(db)
	code := m.Run()
	cleanDatabase(db, rc)
	sqldb.Close()
	os.Exit(code)
}

func newExpect(t *testing.T) *httpexpect.Expect {
	client := &http.Client{
		Transport: fiberRoundTripper{app: app},
	}
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost",
		Client:   client,
		Reporter: httpexpect.NewRequireReporter(t),
	})
}

func migrateSeeds(db *gorm.DB) error {
	projectRoot := config.FindProjectRoot()
	path := filepath.Join(projectRoot ,"config", "users.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var users []dto.CreateUserFromJsonFileDTO

	if err := json.Unmarshal(data, &users); err != nil {
		return  err
	}

	var usersModel []*model.UserModel

	for _, u := range users {
		usersModel = append(usersModel, &model.UserModel{
			AuthId: uuid.New().String(),
			Username: u.Username,
			Role: u.Role,
		})
	}

	if err := db.Create(usersModel).Error; err != nil {
		return  err
	}
	
	return nil
}

func cleanDatabase(db *gorm.DB, rc *redisLib.Client) {
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.UserModel{}, &model.NewsModel{}, &model.BlockModel{}, &model.ImageModel{})
	rc.FlushDB(context.Background())
	logger.ZapLogger.Info("databases were cleaned")
}

func TestMessage(t *testing.T) {
	e := newExpect(t)
	e.GET("/"). 
	Expect().
	Status(200).JSON(). 
	Object()
}