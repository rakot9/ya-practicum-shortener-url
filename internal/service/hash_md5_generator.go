package service

import (
	"crypto/md5"
	"encoding/hex"
	// "fmt"
	// "log/slog"
	"math/rand"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Generate(URL string, regenKey bool, countRandLetter int) (string, error) {
	hash := md5.Sum([]byte(URL))
	Hash := hex.EncodeToString(hash[:])

	if regenKey {
		Hash += randString(countRandLetter + 1)
	}

	return Hash, nil
}
