package auth

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func hashPassword(salt string, password string) string {
	h := md5.New()
	h.Write([]byte(salt + password))
	return hex.EncodeToString(h.Sum(nil))
}

func formatName(name string) string {
	n := strings.TrimSpace(name)
	c := cases.Title(language.Und)
	return c.String(n)
}
