package commands

import (
	"errors"
	"fmt"
	"msnserver/pkg/clients"
	"msnserver/pkg/database"
	"strings"

	"gorm.io/gorm"
)

func HandleFND(db *gorm.DB, c *clients.Client, args string) error {
	args, _, _ = strings.Cut(args, "\r\n")
	tid, args, err := parseTransactionID(args)
	if err != nil {
		return err
	}

	if !c.Session.Authenticated {
		SendError(c, tid, ERR_NOT_LOGGED_IN)
		return nil
	}

	splitArguments := strings.Fields(args)
	if len(splitArguments) != 5 {
		return errors.New("invalid transaction")
	}

	var fname string
	var lname string
	var city string
	var state string
	var country string

	_, fname, found := strings.Cut(splitArguments[0], "fname=")
	if !found {
		return errors.New("invalid transaction")
	}
	_, lname, found = strings.Cut(splitArguments[1], "lname=")
	if !found {
		return errors.New("invalid transaction")
	}
	_, city, found = strings.Cut(splitArguments[2], "city=")
	if !found {
		return errors.New("invalid transaction")
	}
	_, state, found = strings.Cut(splitArguments[3], "state=")
	if !found {
		return errors.New("invalid transaction")
	}
	_, country, found = strings.Cut(splitArguments[4], "country=")
	if !found {
		return errors.New("invalid transaction")
	}

	// First and last name must be specified
	if fname == "*" || lname == "*" {
		SendError(c, tid, ERR_INVALID_PARAMETER)
		return nil
	}

	// Only allow US users to specify city and state
	if country != "US" && (city != "*" || state != "*") {
		return errors.New("invalid transaction")
	}

	var users []database.User

	query := db.Where("first_name ILIKE ? AND last_name ILIKE ?", "%"+fname+"%", "%"+lname+"%")

	if country != "*" {
		query = query.Where("country = ?", country)
	}
	if state != "*" {
		query = query.Where("state = ?", state)
	}
	if city != "*" {
		query = query.Where("city = ?", city)
	}

	query.Limit(100).Find(&users)
	if query.Error != nil {
		return query.Error
	}
	if query.RowsAffected == 100 {
		SendError(c, tid, ERR_TOO_MANY_RESULTS)
		return nil
	}

	if len(users) == 0 {
		res := fmt.Sprintf("FND %s 0 0\r\n", tid)
		c.Send(res)
	}
	for i, user := range users {
		res := fmt.Sprintf("FND %s %d %d fname=%s lname=%s city=%s state=%s country=%s\r\n",
			tid, i+1, len(users), user.FirstName, user.LastName, user.City, user.State, user.Country)
		c.Send(res)
	}

	return nil
}
