package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func FormatString(name string) string {
	n := strings.TrimSpace(name)
	c := cases.Title(language.Und)
	return c.String(n)
}
