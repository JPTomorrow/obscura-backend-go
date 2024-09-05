package auth

import (
	"crypto/rand"
	"encoding/hex"
)

var tokens []string

type BearerToken struct {
	Token string `json:"token"`
}

func GenerateLoginToken() (BearerToken, error) {
	token, err := randomHex(20)
	if err != nil {
		return BearerToken{}, err
	}
	tokens = append(tokens, token)
	i := BearerToken{token}
	return i, nil
}

func IsValidLoginToken(bearerToken string) bool {
	for _, token := range tokens {
		if token == bearerToken {
			return true
		}
	}
	return false
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
