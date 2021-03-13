package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const TrackingRequestHeader = "X-Request-ID"

type tracking struct {
	Logger *zap.Logger
}

func Logger(logger *zap.Logger) *tracking {
	return &tracking{
		Logger: logger,
	}
}

func (l *tracking) LogWithContext(c *fiber.Ctx) error {
	rID := c.Get(TrackingRequestHeader)
	l.Logger.Info("incoming request tracking", zap.String(TrackingRequestHeader, rID))
	return c.Next()
}
