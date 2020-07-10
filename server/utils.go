package server

import (
	"errors"
	"fmt"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

func ValidateJWT(jwtToken string) (bool, error) {
	if len(jwtToken) == 0 {
		return false, errors.New("Missing token")
	}

	_, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return JwtSecret, nil
	})

	if err != nil {
		return false, errors.New("Invalid token")
	}

	return true, nil
}

func GetEnvWithDefault(key string, fallback string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		value = fallback
	}

	return value
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return JwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if token.Valid == false || ok == false {
		return nil, err
	}

	return claims, nil
}
