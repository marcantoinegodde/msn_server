package commands

import (
	"errors"
	"fmt"
	"log"
	"msnserver/pkg/database"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func HandleADD(c chan string, db *gorm.DB, s *Session, args string) error {
	// TODO: Add asynchronous communication reverse list
	// TODO: Add group number to forward list

	args, _, _ = strings.Cut(args, "\r\n")
	transactionID, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	splitArguments := strings.Fields(args)
	if len(splitArguments) < 3 || len(splitArguments) > 4 {
		return errors.New("invalid transaction")
	}

	listName := splitArguments[0]
	email := splitArguments[1]
	displayName := splitArguments[2]
	var groupNum int
	if len(splitArguments) == 4 {
		groupNum, err = strconv.Atoi(splitArguments[3])
		if err != nil {
			return fmt.Errorf("invalid group number: %w", err)
		}
	}

	if !s.connected {
		SendError(c, transactionID, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	if !isValidEmail(email) {
		SendError(c, transactionID, ERR_INVALID_PARAMETER)
		log.Printf("Error: invalid email: %s\n", email)
		return nil
	}

	if s.Email == email {
		SendError(c, transactionID, ERR_INVALID_USER)
		log.Printf("Error: tried to add self to list\n")
		return nil
	}

	var user database.User
	query := db.Preload("ForwardList").Preload("AllowList").Preload("BlockList").First(&user, "email = ?", s.Email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	var principal database.User
	query = db.Preload("ForwardList").Preload("AllowList").Preload("BlockList").First(&principal, "email = ?", email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		SendError(c, transactionID, ERR_INVALID_USER)
		log.Printf("Error: user not found: %s\n", email)
		return nil
	} else if query.Error != nil {
		return query.Error
	}

	switch listName {
	case "FL":
		if len(user.ForwardList) >= 150 {
			SendError(c, transactionID, ERR_LIST_FULL)
			log.Printf("Error: forward list full\n")
			return nil
		}

		if isMember(user.ForwardList, &principal) {
			SendError(c, transactionID, ERR_ALREADY_THERE)
			log.Printf("Error: user already in forward list\n")
			return nil
		}

		user.ForwardList = append(user.ForwardList, &principal)
		user.DataVersion++
		if err := db.Save(&user).Error; err != nil {
			return err
		}

		principal.ReverseList = append(principal.ReverseList, &user)
		principal.DataVersion++
		if err := db.Save(&principal).Error; err != nil {
			return err
		}

	case "AL":
		if len(user.AllowList) >= 150 {
			SendError(c, transactionID, ERR_LIST_FULL)
			log.Printf("Error: allow list full\n")
			return nil
		}

		if isMember(user.AllowList, &principal) {
			SendError(c, transactionID, ERR_ALREADY_THERE)
			log.Printf("Error: user already in allow list\n")
			return nil
		}

		if isMember(user.BlockList, &principal) {
			SendError(c, transactionID, ERR_ALREADY_IN_OPPOSITE_LIST)
			log.Printf("Error: trying to add in AL and BL\n")
			return nil
		}

		user.AllowList = append(user.AllowList, &principal)
		user.DataVersion++
		if err := db.Save(&user).Error; err != nil {
			return err
		}

	case "BL":
		if len(user.BlockList) >= 150 {
			SendError(c, transactionID, ERR_LIST_FULL)
			log.Printf("Error: block list full\n")
			return nil
		}

		if isMember(user.BlockList, &principal) {
			SendError(c, transactionID, ERR_ALREADY_THERE)
			log.Printf("Error: user already in block list\n")
			return nil
		}

		if isMember(user.AllowList, &principal) {
			SendError(c, transactionID, ERR_ALREADY_IN_OPPOSITE_LIST)
			log.Printf("Error: trying to add in AL and BL\n")
			return nil
		}

		user.BlockList = append(user.BlockList, &principal)
		user.DataVersion++
		if err := db.Save(&user).Error; err != nil {
			return err
		}

	case "RL":
		err := errors.New("tried to add to reverse list")
		return err

	default:
		err := errors.New("unsupported list")
		return err
	}

	var res string
	if groupNum != 0 {
		res = fmt.Sprintf("ADD %s %s %d %s %s %d\r\n", transactionID, listName, user.DataVersion, email, displayName, groupNum)
	} else {
		res = fmt.Sprintf("ADD %s %s %d %s %s\r\n", transactionID, listName, user.DataVersion, email, displayName)
	}
	c <- res

	if listName == "FL" {
		HandleSendILN(c, transactionID, principal.Status, principal.Email, principal.Name)
	}

	return nil
}

func isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func isMember(userList []*database.User, principal *database.User) bool {
	for _, u := range userList {
		if u.Email == principal.Email {
			return true
		}
	}

	return false
}
