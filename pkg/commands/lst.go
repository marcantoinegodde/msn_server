package commands

import (
	"errors"
	"fmt"
	"log"
	"msnserver/pkg/database"
	"net"
	"slices"
	"strings"

	"gorm.io/gorm"
)

var listTypes = []string{"FL", "AL", "BL", "RL"}

func HandleLST(conn net.Conn, db *gorm.DB, s *Session, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	transactionID, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	if !slices.Contains(listTypes, args) {
		return errors.New("invalid list")
	}

	if !s.connected {
		SendError(conn, transactionID, ERR_NOT_LOGGED_IN)
		return errors.New("not logged in")
	}

	var user database.User
	query := db.Preload("ForwardList").Preload("AllowList").Preload("BlockList").Preload("ReverseList").First(&user, "email = ?", s.email)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else if query.Error != nil {
		return query.Error
	}

	if err := HandleSendLST(conn, transactionID, args, &user); err != nil {
		return err
	}
	return nil
}

func HandleSendLST(conn net.Conn, tid string, lt string, u *database.User) error {
	switch lt {
	case "FL":
		for i, f := range u.ForwardList {
			res := fmt.Sprintf("LST %s %s %d %d %d %s %s\r\n", tid, lt, u.DataVersion, i+1, len(u.ForwardList), f.Email, f.Name)
			log.Println(">>>", res)
			conn.Write([]byte(res))
		}
		if len(u.ForwardList) == 0 {
			res := fmt.Sprintf("LST %s %s %d %d %d\r\n", tid, lt, u.DataVersion, 0, 0)
			log.Println(">>>", res)
			conn.Write([]byte(res))
		}
		return nil

	case "AL":
		for i, a := range u.AllowList {
			res := fmt.Sprintf("LST %s %s %d %d %d %s %s\r\n", tid, lt, u.DataVersion, i+1, len(u.AllowList), a.Email, a.Name)
			log.Println(">>>", res)
			conn.Write([]byte(res))
		}
		if len(u.AllowList) == 0 {
			res := fmt.Sprintf("LST %s %s %d %d %d\r\n", tid, lt, u.DataVersion, 0, 0)
			log.Println(">>>", res)
			conn.Write([]byte(res))
		}
		return nil

	case "BL":
		for i, b := range u.BlockList {
			res := fmt.Sprintf("LST %s %s %d %d %d %s %s\r\n", tid, lt, u.DataVersion, i+1, len(u.BlockList), b.Email, b.Name)
			log.Println(">>>", res)
			conn.Write([]byte(res))
		}
		if len(u.BlockList) == 0 {
			res := fmt.Sprintf("LST %s %s %d %d %d\r\n", tid, lt, u.DataVersion, 0, 0)
			log.Println(">>>", res)
			conn.Write([]byte(res))
		}
		return nil

	case "RL":
		for i, r := range u.ReverseList {
			res := fmt.Sprintf("LST %s %s %d %d %d %s %s\r\n", tid, lt, u.DataVersion, i+1, len(u.ReverseList), r.Email, r.Name)
			log.Println(">>>", res)
			conn.Write([]byte(res))
		}
		if len(u.ReverseList) == 0 {
			res := fmt.Sprintf("LST %s %s %d %d %d\r\n", tid, lt, u.DataVersion, 0, 0)
			log.Println(">>>", res)
			conn.Write([]byte(res))
		}
		return nil

	default:
		SendError(conn, tid, ERR_INVALID_PARAMETER)
		return errors.New("invalid list type")
	}
}
