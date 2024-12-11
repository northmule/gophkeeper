package util

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

// PasswordHashSha256 хэш пароля
func PasswordHashSha256(password string) string {
	hashAlg := sha256.New()
	hashAlg.Write([]byte(password))
	return fmt.Sprintf("%x", hashAlg.Sum(nil))
}

// PasswordHashSha512 хэш пароля
func PasswordHashSha512(password string) string {
	hashAlg := sha512.New()
	hashAlg.Write([]byte(password))
	return fmt.Sprintf("%x", hashAlg.Sum(nil))
}
