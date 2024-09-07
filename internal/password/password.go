package password

import (
	"bjungle/blockchain-engine/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

func Compare(id, hashedPassword, p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(p))
	if err != nil {
		logger.Warning.Printf("password is invalid: %s, %v", id, err)
		return false
	}
	return true
}

func Encrypt(password string) string {
	bp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error.Printf("generating hash in password: %v", err)
	}
	return string(bp)
}
