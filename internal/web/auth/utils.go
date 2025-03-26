package auth

import (
	"crypto/md5"
	"encoding/hex"
)

func hashPassword(salt string, password string) string {
	h := md5.New()
	h.Write([]byte(salt + password))
	return hex.EncodeToString(h.Sum(nil))
}
