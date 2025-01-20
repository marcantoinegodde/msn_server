package commands

import (
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"

	"gorm.io/gorm"
)

func HandleSendFLN(db *gorm.DB, clients map[string]*clients.Client, s *clients.Session) error {
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

		res := fmt.Sprintf("FLN %s\r\n", user.Email)
		clients[contact.Email].SendChan <- res
	}

	return nil
}
