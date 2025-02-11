package commands

import (
	"errors"
	"fmt"
	"log"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
	"strings"
	"sync"

	"gorm.io/gorm"
)

func HandleREM(db *gorm.DB, m *sync.Mutex, clients map[string]*clients.Client, c *clients.Client, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	tid, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	if !c.Session.Authenticated {
		SendError(c, tid, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	splitArguments := strings.Fields(args)
	if len(splitArguments) != 2 {
		return errors.New("invalid transaction")
	}

	listName := splitArguments[0]
	email := splitArguments[1]

	var user database.User
	query := db.First(&user, "email = ?", c.Session.Email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	switch listName {
	case "FL":
		// Find principal to remove from forward list
		var principal database.User
		err := db.Model(&user).Where("email = ?", email).Association("ForwardList").Find(&principal)
		if err != nil {
			SendError(c, tid, ERR_NOT_ON_LIST)
			log.Println("Error: tried to remove user not on list")
			return nil
		}

		// Remove principal from forward list
		if err := db.Model(&user).Association("ForwardList").Delete(&principal); err != nil {
			return err
		}

		// Update user's data version
		user.DataVersion++
		if err := db.Save(&user).Error; err != nil {
			return err
		}

		// Remove user from principal's reverse list
		if err := db.Model(&principal).Association("ReverseList").Delete(&user); err != nil {
			return err
		}

		// Update principal's data version
		principal.DataVersion++
		if err := db.Save(&principal).Error; err != nil {
			return err
		}

		// Notify principal if online
		m.Lock()
		principalClient, ok := clients[principal.Email]
		if ok {
			res := fmt.Sprintf("REM %s %s %d %s\r\n", "0", "RL", principal.DataVersion, user.Email)
			principalClient.Send(res)
		}
		m.Unlock()

	case "AL":
		// Find principal to remove from allow list
		var principal database.User
		err := db.Model(&user).Where("email = ?", email).Association("AllowList").Find(&principal)
		if err != nil {
			SendError(c, tid, ERR_NOT_ON_LIST)
			log.Println("Error: tried to remove user not on list")
			return nil
		}

		// Remove principal from allow list
		if err := db.Model(&user).Association("AllowList").Delete(&principal); err != nil {
			return err
		}

		// Update user's data version
		user.DataVersion++
		if err := db.Save(&user).Error; err != nil {
			return err
		}

	case "BL":
		// Find principal to remove from block list
		var principal database.User
		err := db.Model(&user).Where("email = ?", email).Association("BlockList").Find(&principal)
		if err != nil {
			SendError(c, tid, ERR_NOT_ON_LIST)
			log.Println("Error: tried to remove user not on list")
			return nil
		}

		// Remove principal from block list
		if err := db.Model(&user).Association("BlockList").Delete(&principal); err != nil {
			return err
		}

		// Update user's data version
		user.DataVersion++
		if err := db.Save(&user).Error; err != nil {
			return err
		}

	case "RL":
		// User cannot modify reverse list
		err := errors.New("tried to add to reverse list")
		return err

	default:
		err := errors.New("unsupported list")
		return err
	}

	res := fmt.Sprintf("REM %d %s %d %s\r\n", tid, listName, user.DataVersion, email)
	c.Send(res)

	return nil
}
