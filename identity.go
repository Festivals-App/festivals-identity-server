package identity

import "github.com/Festivals-App/festivals-identity-server/authentication"

func functiontousetheshit() {
	authentication.IsAuthenticated([]string{"a", "b"}, nil)
}

/*
func homePrint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Super Nice Info")
}

func main() {

	fmt.Println("My server")
	tokenString, error := token.GenerateJWT()

	if error != nil {
		fmt.Println("Error generating token string")
	}

	// authmenow "github.com/Festivals-App/festivals-identity-server/asfojboafsIHB"

	fmt.Println(tokenString)

	http.Handle("/", authentication.IsAuthenticated([]string{"a", "b"}, homePrint))

	log.Fatal(http.ListenAndServe(":9000", nil))


		conf := config.DefaultConfig()
		if len(os.Args) > 1 {
			conf = config.ParseConfig(os.Args[1])
		}

		serverInstance := &server.Server{}
		serverInstance.Initialize(conf)
		serverInstance.Run(":" + strconv.Itoa(conf.ServicePort))

}
*/
