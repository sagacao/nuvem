package utils

import (
	"fmt"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"
)

func GenerateUUID(openid string) string {
	uid, err := uuid.NewV4()
	if err != nil {
		log.Printf("Something went wrong: %s", err)
		return fmt.Sprintf("%s-%v", openid, time.Now().Unix())
	}
	return uid.String()
}
