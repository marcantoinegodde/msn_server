package commands

import (
	"errors"
	"fmt"
	"log"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
	"strings"

	"gorm.io/gorm"
)

func HandleREM(c chan string, db *gorm.DB, s *clients.Session, clients map[string]*clients.Client, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	transactionID, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	splitArguments := strings.Fields(args)
	if len(splitArguments) != 2 {
		return errors.New("invalid transaction")
	}

	listName := splitArguments[0]
	email := splitArguments[1]

	if !s.Authenticated {
		SendError(c, transactionID, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	var user database.User
	query := db.First(&user, "email = ?", s.Email)
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
			SendError(c, transactionID, ERR_NOT_ON_LIST)
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
		if clients[principal.Email] != nil {
			res := fmt.Sprintf("REM %s %s %d %s\r\n", "0", "RL", principal.DataVersion, user.Email)
			clients[principal.Email].SendChan <- res
		}

	case "AL":
		// Find principal to remove from allow list
		var principal database.User
		err := db.Model(&user).Where("email = ?", email).Association("AllowList").Find(&principal)
		if err != nil {
			SendError(c, transactionID, ERR_NOT_ON_LIST)
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
			SendError(c, transactionID, ERR_NOT_ON_LIST)
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

	res := fmt.Sprintf("REM %s %s %d %s\r\n", transactionID, listName, user.DataVersion, email)
	c <- res

	return nil
}
