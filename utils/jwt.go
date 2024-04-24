package utils

import (
	"net/http"
	"organization/config"
	error_handling "organization/error"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	jwt.StandardClaims
}

func CreateJWT(audience string, subject string) (string, error) {
	jwtKey := []byte(config.ConfigVal.JWTKey)
	expirationTime := time.Now().Add(time.Minute * 5)

	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Audience:  audience,
			Subject:   subject,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyJWT(token string, audience string, subject string) error {
	jwtKey := []byte(config.ConfigVal.JWTKey)
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return error_handling.JWTErrSignatureInvalid
		}
		return error_handling.CustomError{StatusCode: http.StatusInternalServerError, ErrorMessage: err.Error()}
	}

	if !tkn.Valid {
		return error_handling.JWTTokenInvalid
	}
	if claims.Audience != audience && claims.Subject == subject {
		return error_handling.JWTTokenInvalidDetails
	}
	return nil
}
