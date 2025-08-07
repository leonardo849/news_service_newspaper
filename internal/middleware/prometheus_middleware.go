package middleware

import (
	"template-backend/internal/prometheus"
	"time"

	"fmt"

	"github.com/gofiber/fiber/v2"
)

func PrometheusMiddleware() fiber.Handler {
    return  func(c *fiber.Ctx) error {
        start := time.Now()
        duration := time.Since(start)

        status := c.Response().StatusCode()
        path := c.Route().Path 
        method := c.Method()
        prometheus.HttpRequestsTotal.WithLabelValues(method, path, fmt.Sprint(status)).Inc()
        prometheus.HttpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())

        return c.Next()
    }
}