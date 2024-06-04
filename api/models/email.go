package models

type Email struct {
	//gorm.Model
	//ID    uint   `json:"id" gorm:"primary_key"`
	Email string `gorm:"primary_key"`
}
