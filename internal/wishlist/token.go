package wishlist

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateAccessToken() string{
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}