package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"
	length   = 10
)

func GenerateShortID() (string, error) {
	result := make([]byte, length)
	max := big.NewInt(int64(len(alphabet)))

	for i := range length {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		result[i] = alphabet[n.Int64()]
	}
	return string(result), nil
}
