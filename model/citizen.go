package model

import "gorm.io/gorm"

type Citizen struct {
	gorm.Model
	CID *string `gorm:"column:cid"`
}
