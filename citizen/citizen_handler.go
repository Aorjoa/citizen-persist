package citizen

import (
	"net/http"
	"time"

	"github.com/Aorjoa/citizen-persist/constant"
	"github.com/Aorjoa/citizen-persist/mq"
	"github.com/Aorjoa/citizen-persist/redis"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// CitizenRequest should use to map incoming request
type CitizenRequest struct {
	CitizenID string `json:"citizen_id"`
}

// ErrorMessageResponse should use to respose message as an application standard
type ErrorMessageResponse struct {
	Message string `json:"message"`
}

type citizen struct {
	Logger  *zap.Logger
	KafkaMQ mq.KafkaMQ
	Redis   redis.DBStorage
}

// NewCitizenHandler should create new citizen handler with dependencies inject parameter
func NewHandler(logger *zap.Logger, mq mq.KafkaMQ, redis redis.DBStorage) *citizen {
	return &citizen{
		Logger:  logger,
		KafkaMQ: mq,
		Redis:   redis,
	}
}

// PutCitizenIDToQueue should insert citizen ID from client request into message queue
// and check whether throttle submittion data
func (ci *citizen) PutCitizenIDToQueue(c *fiber.Ctx) error {
	var cit CitizenRequest
	if err := c.BodyParser(&cit); err != nil {
		msg := "unable to parse request"
		ci.Logger.Error(msg)
		return c.Status(http.StatusBadRequest).JSON(Error(msg))
	}

	_, err := ci.Redis.GetData(cit.CitizenID)
	if err == nil {
		return c.SendStatus(http.StatusConflict)
	}

	err = ci.Redis.SetData(cit.CitizenID, true, 10*time.Second)
	if err != nil {
		msg := "unable to set citizen id to redis"
		ci.Logger.Error(msg)
		return c.Status(http.StatusInternalServerError).JSON(Error(msg))
	}

	ci.KafkaMQ.Push([]byte(constant.CitizenMesssage), []byte(cit.CitizenID))

	return c.SendStatus(http.StatusNoContent)
}

func Error(message string) ErrorMessageResponse {
	return ErrorMessageResponse{Message: message}
}
