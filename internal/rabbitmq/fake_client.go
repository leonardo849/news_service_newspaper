package rabbitmq

import (
	"fmt"
	"news_service/internal/logger"
	"time"
)

type fakeClient struct {
}

func (c *fakeClient) declareExchanges() error {
	
	logger.ZapLogger.Info("[fake]declaring exchanges")
	return  nil
}

func (c *fakeClient) consumeTopicUserAuth()  {
	 logger.ZapLogger.Info("[fake] consuming")

    go func() {
        
        for {
            time.Sleep(2 * time.Second)
            fmt.Println("fake message received")
        }
    }()
	select{}
}

