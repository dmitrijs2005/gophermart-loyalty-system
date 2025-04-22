package auth

import (
	"time"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/common"
	"github.com/golang-jwt/jwt/v4"
)

// Claims — структура утверждений, которая включает стандартные утверждения и
// одно пользовательское UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func GenerateToken(userID string, secretKey string, validityDuration *time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(*validityDuration)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserIDFromToken(tokenString string, secretKey string) (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", common.ErrorInvalidToken
	}

	return claims.UserID, nil
}
