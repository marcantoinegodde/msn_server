package commands

import (
	"errors"
	"fmt"
	"log"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
	"msnserver/pkg/utils"
	"strings"
	"sync"

	"gorm.io/gorm"
)

const (
	MAX_FORWARD_LIST_SIZE int = 150
)

func HandleADD(db *gorm.DB, m *sync.Mutex, clients map[string]*clients.Client, c *clients.Client, args string) error {
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
	if len(splitArguments) != 3 {
		return errors.New("invalid transaction")
	}

	lt := ListType(splitArguments[0])
	email := splitArguments[1]
	displayName := splitArguments[2]

	switch lt {
	case ForwardList, AllowList, BlockList, ReverseList:
		break
	default:
		return errors.New("invalid list")
	}

	if !utils.IsValidEmail(email) {
		SendError(c, tid, ERR_INVALID_PARAMETER)
		log.Printf("Error: invalid email: %s\n", email)
		return nil
	}

	if c.Session.Email == email {
		SendError(c, tid, ERR_INVALID_USER)
		log.Println("Error: tried to add self to list")
		return nil
	}

	var user database.User
	query := db.Preload("ForwardList").Preload("AllowList").Preload("BlockList").First(&user, "email = ?", c.Session.Email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	var principal database.User
	query = db.Preload("ForwardList").Preload("AllowList").Preload("BlockList").First(&principal, "email = ?", email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		SendError(c, tid, ERR_INVALID_USER)
		log.Printf("Error: user not found: %s\n", email)
		return nil
	} else if query.Error != nil {
		return query.Error
	}

	switch lt {
	case ForwardList:
		if len(user.ForwardList) >= MAX_FORWARD_LIST_SIZE {
			SendError(c, tid, ERR_LIST_FULL)
			log.Println("Error: forward list full")
			return nil
		}

		if isMember(user.ForwardList, &principal) {
			SendError(c, tid, ERR_ALREADY_THERE)
			log.Println("Error: user already in forward list")
			return nil
		}

		// Add principal to user's forward list
		user.ForwardList = append(user.ForwardList, &principal)
		user.DataVersion++
		if err := db.Save(&user).Error; err != nil {
			return err
		}

		// Add user to principal's reverse list
		principal.ReverseList = append(principal.ReverseList, &user)
		principal.DataVersion++
		if err := db.Save(&principal).Error; err != nil {
			return err
		}

		// Notify principal if online
		m.Lock()
		principalClient, ok := clients[principal.Email]
		if ok {
			res := fmt.Sprintf("ADD %s %s %d %s %s\r\n", "0", ReverseList, principal.DataVersion, user.Email, user.DisplayName)
			principalClient.Send(res)
		}
		m.Unlock()

		// Notify user if online, not blocked and explicitely allowed if BLP is BL
		if !(principal.Status == database.FLN || principal.Status == database.HDN) &&
			!isMember(principal.BlockList, &user) &&
			!(principal.Blp == database.BL && !isMember(principal.AllowList, &user)) {
			HandleSendILN(c, tid, principal.Status, principal.Email, principal.DisplayName)
		}

	case AllowList:
		if isMember(user.BlockList, &principal) {
			SendError(c, tid, ERR_ALREADY_IN_OPPOSITE_LIST)
			log.Println("Error: trying to add in AL and BL")
			return nil
		}

		// Add principal to user's allow list
		user.AllowList = append(user.AllowList, &principal)
		user.DataVersion++
		if err := db.Save(&user).Error; err != nil {
			return err
		}

		// Notify principal if online
		m.Lock()
		principalClient, ok := clients[principal.Email]
		if ok {
			HandleSendNLN(principalClient, user.Status, user.Email, user.DisplayName)
		}
		m.Unlock()

	case BlockList:
		if isMember(user.AllowList, &principal) {
			SendError(c, tid, ERR_ALREADY_IN_OPPOSITE_LIST)
			log.Println("Error: trying to add in AL and BL")
			return nil
		}

		// Add principal to user's block list
		user.BlockList = append(user.BlockList, &principal)
		user.DataVersion++
		if err := db.Save(&user).Error; err != nil {
			return err
		}

		// Notify principal if online
		m.Lock()
		principalClient, ok := clients[principal.Email]
		if ok {
			HandleSendFLN(principalClient, user.Email)
		}
		m.Unlock()

	case ReverseList:
		// User cannot modify reverse list
		err := errors.New("tried to add to reverse list")
		return err

	default:
		err := errors.New("invalid list")
		return err
	}

	res := fmt.Sprintf("ADD %d %s %d %s %s\r\n", tid, lt, user.DataVersion, email, displayName)
	c.Send(res)

	return nil
}
