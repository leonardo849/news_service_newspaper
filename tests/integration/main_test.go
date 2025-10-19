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
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	constsSl "github.com/leonardo849/shared_library_news_paper/pkg/consts"
	jwtSl "github.com/leonardo849/shared_library_news_paper/pkg/jwt"
	redisLib "github.com/redis/go-redis/v9"
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

type User struct {
	Role string
	Token string
	Id string
	Username string
}

var journalist User
var developer User
var customer User
var ceo User

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
	secret := os.Getenv("SECRETWORDJWT")
	if err := migrateSeeds(db, secret); err != nil {
		logger.ZapLogger.Error(err.Error())
		os.Exit(1)
	}
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

func generateJwt(id string, updatedAt time.Time, role string, secret string) (string, error) {
	jwt, err := jwtSl.GenerateJWT(id, updatedAt, role, secret)
	if err != nil {
		logger.ZapLogger.Error(err.Error())
		return  "", err
	}
	return  jwt, nil
}

func migrateSeeds(db *gorm.DB, secret string) error {
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
		authId := uuid.New().String()

		token, err := generateJwt(authId, time.Now().Add(-5 *time.Minute), u.Role, secret)
		if err != nil {
			logger.ZapLogger.Error(err.Error())
			return  err
		}
		user := User{
			Role: u.Role,
			Token: token,
			Id: authId,
			Username: u.Username,
		}
		if u.Role == constsSl.Ceo {
			ceo = user
		} else if u.Role == constsSl.Customer {
			customer = user
		} else if u.Role == constsSl.Journalist {
			journalist = user
		} else if u.Role == constsSl.Developer {
			developer = user
		}
		usersModel = append(usersModel, &model.UserModel{
			AuthId: authId,
			Username: u.Username,
			Role: u.Role,
		})
	}

	if err := db.Create(usersModel).Error; err != nil {
		return  err
	}
	var count int64

	db.Find(&model.UserModel{}).Count(&count)
	log.Printf("count: %d", count)
	
	return nil
}

func cleanDatabase(db *gorm.DB, rc *redisLib.Client) {
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Exec(`DELETE FROM authors_news`)
    db.Session(&gorm.Session{AllowGlobalUpdate: true}).Exec(`DELETE FROM image_models`)
    db.Session(&gorm.Session{AllowGlobalUpdate: true}).Exec(`DELETE FROM block_models`)
    db.Session(&gorm.Session{AllowGlobalUpdate: true}).Exec(`DELETE FROM news_models`)
    db.Session(&gorm.Session{AllowGlobalUpdate: true}).Exec(`DELETE FROM user_models`)
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