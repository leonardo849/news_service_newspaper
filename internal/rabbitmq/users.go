package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"
	"news_service/internal/logger"
	"news_service/internal/repository"
	"news_service/internal/service"

	constsSl "github.com/leonardo849/shared_library_news_paper/pkg/consts"
	dtoSl "github.com/leonardo849/shared_library_news_paper/pkg/dto"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const authQueue = "auth_queue"
const topicUserAuth = "user.auth.*"

func (c *client) declareExchanges() error {
	err := c.ch.ExchangeDeclare(
		constsSl.ExchangeNameAuthEvents,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.ZapLogger.Fatal("error declaring exchange", zap.Error(err))
		return  err
	}
	logger.ZapLogger.Info("exchange was declared")
	return  nil
}



func (c *client) createUserFromAuth(input dtoSl.AuthPublishUserCreated, userService *service.UserService) error {
	
	status, message := userService.CreateUser(input)
	if status >= 400 {
		logger.ZapLogger.Warn("error creating verified user auth_id:" + input.AuthId)
		return errors.New(message.(string))
	}
	return  nil
		
	
}

func (c *client) createUsersFromAuth(input []dtoSl.AuthPublishUserCreated, userService *service.UserService) error {
	// status, message := userService.CreateUsers(input)
	// if status >= 400 {
	// 	if status == 400 {
	// 		logger.ZapLogger.Warn("there isn't any valid auth_id")
	// 		return  nil
	// 	} else {
	// 		logger.ZapLogger.Warn("error", zap.Error(errors.New(message)))
	// 		return  errors.New(message)
	// 	}
	// }
	// return  nil
	for _, u := range input {
		status, _ := userService.CreateUser(u)
		if status >= 400 {
			logger.ZapLogger.Warn("error creating verified user auth_id:" + u.AuthId)
			continue
		}
	}
	return  nil
}



func (c *client) consumeTopicUserAuth()  {
	q, err := c.ch.QueueDeclare(authQueue, true, false, false, false, amqp091.Table{
		"x-message-ttl":  int32(10 * 60 * 1000),
	})
	if err != nil {
		logger.ZapLogger.Fatal("error declaring exchange", zap.Error(err))
	}
	err = c.ch.QueueBind(q.Name, topicUserAuth, constsSl.ExchangeNameAuthEvents, false, nil)
	if err != nil {
		logger.ZapLogger.Fatal("error declaring exchange", zap.Error(err))
	}
	msgs, err := c.ch.Consume(
		authQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.ZapLogger.Error("error consuming channel", zap.Error(err))
	}
	
	userRepository := repository.CreateUserRepository(repository.DB)
	userService := service.CreateNewUserService(userRepository)


	go func() {
		
		userRepository.SetDatabase(repository.DB)
		for d :=  range msgs {
			logger.ZapLogger.Info(fmt.Sprintf("one more message coming from auth exchange. Routing key: %s", d.RoutingKey))
			if d.RoutingKey == constsSl.KeyUserAuthVerified {
				var userAuthVerified dtoSl.AuthPublishUserCreated
				err := json.Unmarshal(d.Body, &userAuthVerified)
				if err != nil {
					logger.ZapLogger.Warn("error in json unmarshal", zap.Error(err))
					continue
				}
				if err := c.createUserFromAuth(userAuthVerified, userService); err != nil {
					logger.ZapLogger.Warn("error in create users from auth", zap.Error(err))
					d.Nack(false, true)
					continue
				} else {
					d.Ack(false)
				}
			} else if d.RoutingKey == constsSl.KeyUsersSeed {
				var usersAuthVerified []dtoSl.AuthPublishUserCreated
				err := json.Unmarshal(d.Body, &usersAuthVerified)
				if err != nil {
					logger.ZapLogger.Warn("error in json unmarshal", zap.Error(err))
					continue
				}
				c.createUsersFromAuth(usersAuthVerified, userService)
				d.Ack(false)
			}
		}
	}()
	select{}
}