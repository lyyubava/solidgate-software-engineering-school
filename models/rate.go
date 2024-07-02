package models

import (
	"time"
)

type Rate struct {
	ID           uint `json:"id" gorm:"primary_key"`
	Rate         float32
	ExchangeDate time.Time
}
