package uniq

import (
	uuid "github.com/satori/go.uuid"
	"strings"
)

func GenUniqId() string {
	return strings.ReplaceAll("-", "", uuid.NewV4().String())
}
