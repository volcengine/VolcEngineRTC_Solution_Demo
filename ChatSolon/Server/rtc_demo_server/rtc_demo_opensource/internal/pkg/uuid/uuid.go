package uuid

import (
	u "github.com/satori/go.uuid"
)

func GetUUID() string {
	return u.Must(u.NewV4(), nil).String()
}
