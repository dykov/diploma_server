package util

import (
	"crypto/sha256"
	"fmt"
)

func HashPassword(password, salt1 string) string {

	password = password + salt1
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum(nil))

}
