package citizen_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Aorjoa/citizen-persist/citizen"
	"github.com/Aorjoa/citizen-persist/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/utils"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestPutCitizenIDToQueue_Success_WithHitRedis(t *testing.T) {
	core, _ := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)
	defer logger.Sync()

	app := fiber.New()

	rs := &mockRedis{}
	rs.On("GetData", mock.Anything).Return(true, nil)

	c := citizen.NewHandler(logger, rs)

	rd := strings.NewReader(`{"citizen_id": "1234567890123"}`)

	app.Post("/api/v1/citizens", c.PutCitizenIDToQueue)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/citizens", rd)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	utils.AssertEqual(t, http.StatusConflict, resp.StatusCode, "Status code")
}

func TestPutCitizenIDToQueue_Success_WithoutRedisHit(t *testing.T) {
	core, _ := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)
	defer logger.Sync()

	app := fiber.New()

	rs := &mockRedis{}
	mErr := errors.New("test error")
	rs.On("GetData", mock.Anything).Return(false, mErr)
	rs.On("SetData", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	c := citizen.NewHandler(logger, rs)

	rd := strings.NewReader(`{"citizen_id": "1234567890123"}`)

	app.Post("/api/v1/citizens", c.PutCitizenIDToQueue)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/citizens", rd)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "Status code")
}

var _ redis.DBStorage = &mockRedis{}

type mockRedis struct {
	mock.Mock
}

func (m *mockRedis) GetData(key string) (interface{}, error) {
	args := m.Mock.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m *mockRedis) SetData(key string, value interface{}, expiration time.Duration) error {
	args := m.Mock.Called(key, value, expiration)
	return args.Error(0)
}
