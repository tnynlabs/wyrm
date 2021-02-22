package utils

import (
	"crypto/rand"
	"encoding/base64"
)

//GenString Generates random string of size n bytes
func GenString(n int) string {
	randBytes := make([]byte, n)
	rand.Read(randBytes)
	return base64.StdEncoding.EncodeToString(randBytes)
}
