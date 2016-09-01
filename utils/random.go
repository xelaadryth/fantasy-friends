package utils

import (
	"crypto/rand"
	"encoding/base64"
)

//GenerateString creates a random string encoded in base64 that has a length that is a multiple of 4 and shorter than
//maxLength
func GenerateString(maxLength int) (string, error) {
	length := maxLength / 4 * 3
	byteString := make([]byte, length)
	_, err := rand.Read(byteString)

	if err != nil {
		return "", err
	}

	encodedString := base64.StdEncoding.EncodeToString(byteString)

	return encodedString, nil
}
