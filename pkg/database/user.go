package database

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Email       string `gorm:"uniqueIndex"`
	Salt        string
	Password    string
	Name        string
	Status      string
	DataVersion uint32 `gorm:"default:0"`
	Gtc         string
	Blp         string
	ForwardList []*User `gorm:"many2many:forward_list"`
	AllowList   []*User `gorm:"many2many:allow_list"`
	BlockList   []*User `gorm:"many2many:block_list"`
	ReverseList []*User `gorm:"many2many:reverse_list"`
}
