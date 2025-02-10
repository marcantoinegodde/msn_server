package commands

import (
	"fmt"
)

const (
	ERR_SYNTAX_ERROR             int = 200
	ERR_INVALID_PARAMETER        int = 201
	ERR_INVALID_USER             int = 205
	ERR_FQDN_MISSING             int = 206
	ERR_ALREADY_LOGIN            int = 207
	ERR_INVALID_USERNAME         int = 208
	ERR_INVALID_FRIENDLY_NAME    int = 209
	ERR_LIST_FULL                int = 210
	ERR_ALREADY_THERE            int = 215
	ERR_NOT_ON_LIST              int = 216
	ERR_ALREADY_IN_THE_MODE      int = 218
	ERR_ALREADY_IN_OPPOSITE_LIST int = 219
	ERR_NO_SUCH_GROUP            int = 231
	ERR_SWITCHBOARD_FAILED       int = 280
	ERR_NOTIFY_XFR_FAILED        int = 281

	ERR_REQUIRED_FIELDS_MISSING int = 300
	ERR_TOO_MANY_RESULTS        int = 301
	ERR_NOT_LOGGED_IN           int = 302

	ERR_INTERNAL_SERVER int = 500
	ERR_DB_SERVER       int = 501
	ERR_FILE_OPERATION  int = 510
	ERR_MEMORY_ALLOC    int = 520
	ERR_WRONG_CHL       int = 540

	ERR_SERVER_BUSY        int = 600
	ERR_SERVER_UNAVAILABLE int = 601
	ERR_PEER_NS_DOWN       int = 602
	ERR_DB_CONNECT         int = 603
	ERR_SERVER_GOING_DOWN  int = 604

	ERR_CREATE_CONNECTION          int = 707
	ERR_CVR_UNKNOWN_OR_NOT_ALLOWED int = 710
	ERR_BLOCKING_WRITE             int = 711
	ERR_SESSION_OVERLOAD           int = 712
	ERR_USER_TOO_ACTIVE            int = 713
	ERR_TOO_MANY_SESSIONS          int = 714
	ERR_NOT_EXPECTED               int = 715
	ERR_BAD_FRIEND_FILE            int = 717

	ERR_AUTHENTICATION_FAILED    int = 911
	ERR_NOT_ALLOWED_WHEN_OFFLINE int = 913
	ERR_NOT_ACCEPTING_NEW_USERS  int = 920
	ERR_PASSPORT_NOT_VERIFIED    int = 924
)

func SendError(c chan string, transactionID string, errorCode int) {
	res := fmt.Sprintf("%d %s\r\n", errorCode, transactionID)
	c <- res
}
