package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func MakeToken(login string) string {
	date := time.Now()
	num := rand.Int()
	hash := sha256.New()
	hash.Write([]byte(date.String()))
	hash.Write([]byte(strconv.Itoa(num)))
	hash.Write([]byte(login))
	hashSum := hash.Sum(nil)
	return hex.EncodeToString(hashSum)
}
