package main

import (
	"os"
	"strconv"
	"time"

	"github.com/Festivals-App/festivals-gateway/server/heartbeat"
	"github.com/Festivals-App/festivals-gateway/server/logger"
	"github.com/Festivals-App/festivals-identity-server/server"
	"github.com/Festivals-App/festivals-identity-server/server/config"
	"github.com/rs/zerolog/log"
)

func main() {

	logger.Initialize("/var/log/festivals-identity-server/info.log", true)

	log.Info().Msg("Server startup.")

	conf := config.DefaultConfig()
	if len(os.Args) > 1 {
		conf = config.ParseConfig(os.Args[1])
	}

	log.Info().Msg("Server configuration was initialized.")

	serverInstance := &server.Server{}
	serverInstance.Initialize(conf)

	go serverInstance.Run(conf.ServiceBindAddress + ":" + strconv.Itoa(conf.ServicePort))
	log.Info().Msg("Server did start.")

	go sendHeartbeat(conf)
	log.Info().Msg("Heartbeat routine was started.")

	// wait forever
	// https://stackoverflow.com/questions/36419054/go-projects-main-goroutine-sleep-forever
	select {}
}

func sendHeartbeat(conf *config.Config) {
	for {
		timer := time.After(time.Second * 2)
		<-timer
		var beat *heartbeat.Heartbeat = &heartbeat.Heartbeat{Service: "festivals-identity-server", Host: conf.ServiceBindAddress, Port: conf.ServicePort, Available: true}
		err := heartbeat.SendHeartbeat(conf.LoversEar, conf.ServiceKey, beat)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send heartbeat")
		}
	}
}
