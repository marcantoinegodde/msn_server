package commands

import (
	"fmt"
	"log"
	"net"
)

const (
	ERR_SYNTAX_ERROR             = 200
	ERR_INVALID_PARAMETER        = 201
	ERR_INVALID_USER             = 205
	ERR_FQDN_MISSING             = 206
	ERR_ALREADY_LOGIN            = 207
	ERR_INVALID_USERNAME         = 208
	ERR_INVALID_FRIENDLY_NAME    = 209
	ERR_LIST_FULL                = 210
	ERR_ALREADY_THERE            = 215
	ERR_NOT_ON_LIST              = 216
	ERR_ALREADY_IN_THE_MODE      = 218
	ERR_ALREADY_IN_OPPOSITE_LIST = 219
	ERR_NO_SUCH_GROUP            = 231
	ERR_SWITCHBOARD_FAILED       = 280
	ERR_NOTIFY_XFR_FAILED        = 281

	ERR_REQUIRED_FIELDS_MISSING = 300
	ERR_NOT_LOGGED_IN           = 302

	ERR_INTERNAL_SERVER = 500
	ERR_DB_SERVER       = 501
	ERR_FILE_OPERATION  = 510
	ERR_MEMORY_ALLOC    = 520
	ERR_WRONG_CHL       = 540

	ERR_SERVER_BUSY        = 600
	ERR_SERVER_UNAVAILABLE = 601
	ERR_PEER_NS_DOWN       = 602
	ERR_DB_CONNECT         = 603
	ERR_SERVER_GOING_DOWN  = 604

	ERR_CREATE_CONNECTION          = 707
	ERR_CVR_UNKNOWN_OR_NOT_ALLOWED = 710
	ERR_BLOCKING_WRITE             = 711
	ERR_SESSION_OVERLOAD           = 712
	ERR_USER_TOO_ACTIVE            = 713
	ERR_TOO_MANY_SESSIONS          = 714
	ERR_NOT_EXPECTED               = 715
	ERR_BAD_FRIEND_FILE            = 717

	ERR_AUTHENTICATION_FAILED    = 911
	ERR_NOT_ALLOWED_WHEN_OFFLINE = 913
	ERR_NOT_ACCEPTING_NEW_USERS  = 920
	ERR_PASSPORT_NOT_VERIFIED    = 924
)

func SendError(conn net.Conn, transactionID string, errorCode int) {
	res := fmt.Sprintf("%d %s\r\n", errorCode, transactionID)
	log.Println(">>>", res)
	conn.Write([]byte(res))
}
