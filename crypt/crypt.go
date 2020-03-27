package crypt

import (
	"crypto/sha256"
)

func Sha256Sum(input string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	sha := hasher.Sum(nil)
	return sha
}
