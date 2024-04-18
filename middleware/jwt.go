package middleware

import (
	"organization/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey []byte

type Claims struct {
	jwt.StandardClaims
}

func CreateJWT() (string, error) {
	jwtKey = []byte(config.ConfigVal.JWTKey)
	expirationTime := time.Now().Add(time.Minute * 5)

	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Audience:  "User",
			Subject:   "Member details of organization",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
