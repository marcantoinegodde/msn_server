package commands

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

func parseTransactionID(arguments string) (string, string, error) {
	transactionID, arguments, _ := strings.Cut(arguments, " ")

	parsedTransactionID, err := strconv.Atoi(transactionID)
	if err != nil {
		return "", "", errors.New("Invalid transaction ID")
	}

	if parsedTransactionID < 0 {
		return "", "", errors.New("Invalid transaction ID")
	}
	if parsedTransactionID > (int(math.Pow(2, 32)) - 1) {
		return "", "", errors.New("Invalid transaction ID")
	}

	return transactionID, arguments, nil
}
