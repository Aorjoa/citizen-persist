package citizen

import (
	"github.com/Aorjoa/citizen-persist/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type persistent struct {
	Logger *zap.Logger
	DB     *gorm.DB
}

// NewPersistent should save data from message queue to db
func NewPersistent(logger *zap.Logger, db *gorm.DB) *persistent {
	return &persistent{
		Logger: logger,
		DB:     db,
	}
}

func (p *persistent) Create(ct *model.Citizen) error {
	return p.DB.Create(ct).Error
}
