package main

import (
	"net/http"
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

	logger.InitializeGlobalLogger("/var/log/festivals-identity-server/info.log", true)

	log.Info().Msg("Server startup.")

	conf := config.DefaultConfig()
	if len(os.Args) > 1 {
		conf = config.ParseConfig(os.Args[1])
	}

	log.Info().Msg("Server configuration was initialized.")

	serverInstance := &server.Server{}
	serverInstance.Initialize(conf)

	// Redirect traffic from port 80 to used port
	go http.ListenAndServe(":80", serverInstance.CertManager.HTTPHandler(nil))
	log.Info().Msg("Start redirecting http traffic from port 80 to port " + strconv.Itoa(serverInstance.Config.ServicePort) + " and https")

	go serverInstance.Run(conf)
	log.Info().Msg("Server did start.")

	go sendHeartbeat(conf)
	log.Info().Msg("Heartbeat routine was started.")

	// wait forever
	// https://stackoverflow.com/questions/36419054/go-projects-main-goroutine-sleep-forever
	select {}
}

func sendHeartbeat(conf *config.Config) {

	heartbeatClient, err := heartbeat.HeartbeatClient(conf.TLSCert, conf.TLSKey)
	if err != nil {
		log.Fatal().Err(err).Str("type", "server").Msg("Failed to create heartbeat client")
	}
	var beat *heartbeat.Heartbeat = &heartbeat.Heartbeat{Service: "festivals-identity-server", Host: "https://" + conf.ServiceBindHost, Port: conf.ServicePort, Available: true}

	t := time.NewTicker(time.Duration(conf.Interval) * time.Second)
	defer t.Stop()
	for range t.C {
		err = heartbeat.SendHeartbeat(heartbeatClient, conf.LoversEar, conf.ServiceKey, beat)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send heartbeat")
		}
	}
}
