package models

import "gorm.io/gorm"

type Email struct {
	gorm.Model
	Email string
}
