package util

import (
	uuid "github.com/satori/go.uuid"
)

func CreateUUid() string {
	u1 := uuid.NewV4()
	return u1.String()
}
