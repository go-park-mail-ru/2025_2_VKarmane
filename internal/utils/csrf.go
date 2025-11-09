package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ClaimsCSRF struct {
	CSRFID    string    `json:"csrf_id"`
	TimeStamp time.Time `json:"timestamp"`
	jwt.RegisteredClaims
}

func GenerateCSRF(timeStamp time.Time, secret string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	claims := ClaimsCSRF{
		CSRFID:    base64.RawStdEncoding.EncodeToString(b),
		TimeStamp: timeStamp,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateCSRF(tokenString string, secret string) (*ClaimsCSRF, error) {
	claims := &ClaimsCSRF{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	return claims, nil
}
