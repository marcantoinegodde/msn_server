package commands

import (
	"fmt"
)

func HandleSendILN(c chan string, transactionID string, status string, email string, name string) {
	res := fmt.Sprintf("ILN %s %s %s %s\r\n", transactionID, status, email, name)
	c <- res
}
