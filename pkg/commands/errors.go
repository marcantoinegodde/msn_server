package commands

import (
	"fmt"
	"log"
	"net"
)

const (
	ERR_SYNTAX_ERROR                   = 200
	ERR_INVALID_PARAMETER              = 201
	ERR_INVALID_USER                   = 205
	ERR_DOMAIN_NAME_MISSING            = 206
	ERR_ALREADY_LOGGED_IN              = 207
	ERR_INVALID_USERNAME               = 208
	ERR_INVALID_FUSERNAME              = 209
	ERR_USER_LIST_FULL                 = 210
	ERR_USER_ALREADY_THERE             = 215
	ERR_USER_ALREADY_ON_LIST           = 216
	ERR_USER_NOT_ONLINE                = 217
	ERR_ALREADY_IN_MODE                = 218
	ERR_USER_IN_OPPOSITE_LIST          = 219
	ERR_USER_IN_OPPOSITE_LIST2         = 220
	ERR_ADD_CONTACT_TO_NON_EXISTENT    = 231
	ERR_SWITCHBOARD_FAILED             = 280
	ERR_TRANSFER_TO_SWITCHBOARD_FAILED = 281

	ERR_AUTH_FAILED                   = 911
	ERR_NOT_ALLOWED_WHEN_OFFLINE      = 913
	ERR_NOT_ACCEPTING_NEW_USERS       = 920
	ERR_PASSPORT_ACCOUNT_NOT_VERIFIED = 924
)

func SendError(conn net.Conn, transactionID string, errorCode int) {
	res := fmt.Sprintf("%d %s\r\n", errorCode, transactionID)
	log.Println(">>>", res)
	conn.Write([]byte(res))
}
