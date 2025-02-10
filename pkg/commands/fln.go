package commands

import (
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
	"sync"

	"gorm.io/gorm"
)

func HandleBatchFLN(db *gorm.DB, m *sync.Mutex, clients map[string]*clients.Client, c *clients.Client) error {
	var user database.User
	query := db.Preload("AllowList").Preload("BlockList").Preload("ReverseList").First(&user, "email = ?", c.Session.Email)
	if query.Error != nil {
		return query.Error
	}

	for _, contact := range user.ReverseList {
		if contact.Status == "FLN" {
			continue
		}

		if isMember(user.BlockList, contact) {
			continue
		}

		if user.Blp == "BL" && !isMember(user.AllowList, contact) {
			continue
		}

		m.Lock()
		contactClient, ok := clients[contact.Email]
		if !ok {
			continue
		}

		HandleSendFLN(contactClient, user.Email)
		m.Unlock()
	}

	return nil
}

func HandleSendFLN(c *clients.Client, email string) {
	res := fmt.Sprintf("FLN %s\r\n", email)
	c.Send(res)
}
