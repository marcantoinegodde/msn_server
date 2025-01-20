package commands

import (
	"errors"
	"fmt"
	"log"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
	"slices"
	"strings"

	"gorm.io/gorm"
)

var statusCodes = []string{"NLN", "FLN", "HDN", "IDL", "AWY", "BSY", "BRB", "PHN", "LUN"}

func HandleCHG(c chan string, db *gorm.DB, s *clients.Session, clients map[string]*clients.Client, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	transactionID, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}
	args, _, _ = strings.Cut(args, " ") // Remove the trailing space sent for this command

	if !slices.Contains(statusCodes, args) {
		return fmt.Errorf("invalid status code: %s", args)
	}

	if !s.Authenticated {
		SendError(c, transactionID, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	// Perform nested preloading to load users lists of contacts on user's forward list
	var user database.User
	query := db.Preload("ForwardList.ForwardList").Preload("ForwardList.AllowList").Preload("ForwardList.BlockList").Preload("ReverseList").First(&user, "email = ?", s.Email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	user.Status = args
	query = db.Save(&user)
	if query.Error != nil {
		return query.Error
	}

	res := fmt.Sprintf("CHG %s %s\r\n", transactionID, user.Status)
	c <- res

	// Receive ILN on first CHG
	if !s.InitialPresenceNotification {
		s.InitialPresenceNotification = true

		for _, contact := range user.ForwardList {
			// Skip contacts that are offline or hidden
			if contact.Status == "FLN" || contact.Status == "HDN" {
				continue
			}

			// Skip contacts that have the user on their block list
			if isMember(contact.BlockList, &user) {
				continue
			}

			// Skip contacts in BL mode that don't have the user on their allow list
			if contact.Blp == "BL" && !isMember(contact.AllowList, &user) {
				continue
			}

			// Send initial presence notification
			HandleSendILN(c, transactionID, contact.Status, contact.Email, contact.Name)
		}
	}

	// Inform followers (RL) of the status change
	if user.Status == "HDN" {
		if err := HandleSendFLN(db, clients, s); err != nil {
			log.Println("Error:", err)
		}
	} else {
		if err := HandleSendNLN(db, clients, s); err != nil {
			log.Println("Error:", err)
		}
	}

	return nil
}
