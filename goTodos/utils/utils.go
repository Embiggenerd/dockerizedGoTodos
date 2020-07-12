package utils

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// RandHex simply creates random hex string
// To use to query session data
func RandHex(n int) (string, *HTTPError) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", &HTTPError{
			Code: http.StatusInternalServerError,
			Msg:  "An error happend that isn't your fault",
		}
	}
	return hex.EncodeToString(bytes), nil
}
