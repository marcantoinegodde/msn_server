package commands

import (
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"

	"gorm.io/gorm"
)

func HandleBatchFLN(db *gorm.DB, clients map[string]*clients.Client, s *clients.Session) error {
	var user database.User
	query := db.Preload("AllowList").Preload("BlockList").Preload("ReverseList").First(&user, "email = ?", s.Email)
	if query.Error != nil {
		return query.Error
	}

	for _, contact := range user.ReverseList {
		if contact.Status == "FLN" {
			continue
		}

		if clients[contact.Email] == nil {
			continue
		}

		if isMember(user.BlockList, contact) {
			continue
		}

		if user.Blp == "BL" && !isMember(user.AllowList, contact) {
			continue
		}

		HandleSendFLN(clients[contact.Email].SendChan, user.Email)
	}

	return nil
}

func HandleSendFLN(c chan string, email string) {
	res := fmt.Sprintf("FLN %s\r\n", email)
	c <- res
}
