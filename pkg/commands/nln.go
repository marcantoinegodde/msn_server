package commands

import (
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"

	"gorm.io/gorm"
)

func HandleBatchNLN(db *gorm.DB, clients map[string]*clients.Client, s *clients.Session) error {
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

		HandleSendNLN(clients[contact.Email].SendChan, user.Status, user.Email, user.DisplayName)
	}

	return nil
}

func HandleSendNLN(c chan string, status string, email string, name string) {
	res := fmt.Sprintf("NLN %s %s %s\r\n", status, email, name)
	c <- res
}
