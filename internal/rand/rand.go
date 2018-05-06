package rand

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
)

// bytes returns securely generated random bytes.
func bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// HexEncodedString returns the hexadecimal encoding of random bytes.
func HexEncodedString(n int) (string, error) {
	b, err := bytes(n)
	return hex.EncodeToString(b), err
}

// Base64EncodedString returns the base64 encoding of random bytes.
func Base64EncodedString(n int) (string, error) {
	b, err := bytes(n)
	return base64.StdEncoding.EncodeToString(b), err
}
