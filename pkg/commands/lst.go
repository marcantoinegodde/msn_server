package commands

import (
	"errors"
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
	"slices"
	"strings"

	"gorm.io/gorm"
)

var listTypes = []string{"FL", "AL", "BL", "RL"}

func HandleLST(db *gorm.DB, c *clients.Client, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	transactionID, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	if !c.Session.Authenticated {
		SendError(c, transactionID, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	if !slices.Contains(listTypes, args) {
		return errors.New("invalid list")
	}

	var user database.User
	query := db.Preload("ForwardList").Preload("AllowList").Preload("BlockList").Preload("ReverseList").First(&user, "email = ?", c.Session.Email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	if err := HandleSendLST(c, transactionID, args, &user); err != nil {
		return err
	}
	return nil
}

func HandleSendLST(c *clients.Client, tid string, lt string, u *database.User) error {
	switch lt {
	case "FL":
		for i, f := range u.ForwardList {
			res := fmt.Sprintf("LST %s %s %d %d %d %s %s\r\n", tid, lt, u.DataVersion, i+1, len(u.ForwardList), f.Email, f.DisplayName)
			c.Send(res)
		}
		if len(u.ForwardList) == 0 {
			res := fmt.Sprintf("LST %s %s %d %d %d\r\n", tid, lt, u.DataVersion, 0, 0)
			c.Send(res)
		}
		return nil

	case "AL":
		for i, a := range u.AllowList {
			res := fmt.Sprintf("LST %s %s %d %d %d %s %s\r\n", tid, lt, u.DataVersion, i+1, len(u.AllowList), a.Email, a.DisplayName)
			c.Send(res)
		}
		if len(u.AllowList) == 0 {
			res := fmt.Sprintf("LST %s %s %d %d %d\r\n", tid, lt, u.DataVersion, 0, 0)
			c.Send(res)
		}
		return nil

	case "BL":
		for i, b := range u.BlockList {
			res := fmt.Sprintf("LST %s %s %d %d %d %s %s\r\n", tid, lt, u.DataVersion, i+1, len(u.BlockList), b.Email, b.DisplayName)
			c.Send(res)
		}
		if len(u.BlockList) == 0 {
			res := fmt.Sprintf("LST %s %s %d %d %d\r\n", tid, lt, u.DataVersion, 0, 0)
			c.Send(res)
		}
		return nil

	case "RL":
		for i, r := range u.ReverseList {
			res := fmt.Sprintf("LST %s %s %d %d %d %s %s\r\n", tid, lt, u.DataVersion, i+1, len(u.ReverseList), r.Email, r.DisplayName)
			c.Send(res)
		}
		if len(u.ReverseList) == 0 {
			res := fmt.Sprintf("LST %s %s %d %d %d\r\n", tid, lt, u.DataVersion, 0, 0)
			c.Send(res)
		}
		return nil

	default:
		SendError(c, tid, ERR_INVALID_PARAMETER)
		return errors.New("invalid list type")
	}
}
