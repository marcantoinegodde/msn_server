package database

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Email    string `gorm:"uniqueIndex"`
	Salt     string
	Password string
	Name     string
	Status   string
}
