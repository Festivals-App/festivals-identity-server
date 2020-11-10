package main

import (
	"github.com/Festivals-App/festivals-identity-server/server"
	"github.com/Festivals-App/festivals-identity-server/server/config"
	"os"
	"strconv"
)

func main() {
	conf := config.DefaultConfig()
	if len(os.Args) > 1 {
		conf = config.ParseConfig(os.Args[1])
	}

	serverInstance := &server.Server{}
	serverInstance.Initialize(conf)
	serverInstance.Run(":" + strconv.Itoa(conf.ServicePort))
}
