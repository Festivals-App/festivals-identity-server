package main

import (
	"fmt"
	"log"
	"net/http"

	jwtprocessing "github.com/Festivals-App/festivals-identity-server/jwt"
)

func homePrint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Super Nice Info")
}

func main() {

	fmt.Println("My server")
	tokenString, error := jwtprocessing.GenerateJWT()

	if error != nil {
		fmt.Println("Error generating token string")
	}

	// authmenow "github.com/Festivals-App/festivals-identity-server/authentication"

	fmt.Println(tokenString)

	http.Handle("/", fuckYYOuFuckingGoSchitLanguage([]string{"a", "b"}, homePrint))

	log.Fatal(http.ListenAndServe(":9000", nil))

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
