package model

import "github.com/jinzhu/gorm"

type Example struct {
	gorm.Model
	Name     string
	Phone    string `gorm:"type:char(11);`
	Password string
}
