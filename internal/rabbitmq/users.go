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

func (c *client) createUserFromAuth(input []dtoSl.AuthPublishUserCreated, userService *service.UserService) error {
	if len(input) > 1 {
		status, message := userService.CreateUsers(input)
		if  status >= 400 {
		    logger.ZapLogger.Warn(
				message,
			)
			return errors.New(message)
		} 
	} else {
			status, message := userService.CreateUser(input[0])
			if status >= 400 {
				logger.ZapLogger.Warn(
				message.(string),
			)
			return errors.New(message.(string))
			}
		}
		return  nil
}

func (c *client) consumeTopicUserAuth()  {
	q, err := c.ch.QueueDeclare(authQueue, true, false, false, false, nil)
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
	
	userRepository := repository.CreateUserRepository(repository.DB)
	userService := service.CreateNewUserService(userRepository)


	go func() {
		var usersAuthVerified []dtoSl.AuthPublishUserCreated
		userRepository.SetDatabase(repository.DB)
		for d :=  range msgs {
			logger.ZapLogger.Info(fmt.Sprintf("one more message coming from auth exchange. Routing key: %s", d.RoutingKey))
			if d.RoutingKey == constsSl.KeyUserAuthVerified {
				err = json.Unmarshal(d.Body, &usersAuthVerified)
				if err != nil {
					logger.ZapLogger.Warn("error in json unmarshal", zap.Error(err))
					continue
				}
				if err = c.createUserFromAuth(usersAuthVerified, userService); err != nil {
					logger.ZapLogger.Warn("error in create users from auth", zap.Error(err))
					d.Nack(false, true)
					continue
				} else {
					d.Ack(false)
				}
			}
		}
	}()
	select{}
}