package token

import (
	"fmt"

	"time"

	"github.com/golang-jwt/jwt/v5"
)

var mySigningKey = []byte("mysupersecret")

func GenerateJWT() (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = "Simon Gaus"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}
