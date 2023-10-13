package db

import (
	"encoding/base64"
	"errors"
)

func EncodeNationalID(nationalID string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(nationalID))
	return encoded
}

func DecodeNationalID(encodedNationalID string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encodedNationalID)
	if err != nil {
		return "", errors.New("failed to decode Base64")
	}
	nationalID := string(decoded)
	return nationalID, nil
}
