package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
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

func main() {

	fmt.Println("My server")

	tokenString, error := GenerateJWT()

	if error != nil {
		fmt.Println("Error generating token string")
	}

	fmt.Println(tokenString)
	/*
		conf := config.DefaultConfig()
		if len(os.Args) > 1 {
			conf = config.ParseConfig(os.Args[1])
		}

		serverInstance := &server.Server{}
		serverInstance.Initialize(conf)
		serverInstance.Run(":" + strconv.Itoa(conf.ServicePort))
	*/
}
