package main

import (
	"flag"
	"github.com/opendata-heilbronn/ttn-gateway-exporter/internal/config"
	"github.com/opendata-heilbronn/ttn-gateway-exporter/internal/exporter"
	"github.com/opendata-heilbronn/ttn-gateway-exporter/internal/logging"
	"github.com/opendata-heilbronn/ttn-gateway-exporter/internal/server"
	"github.com/prometheus/client_golang/prometheus"
)

var log = logging.Logger("main")

func main() {
	address := flag.String("address", ":8080", "HTTP listener address")
	targetConfigPath := flag.String("target-config-path", "/etc/ttn-exporter/targets.yaml", "Path to a target config file")
	flag.Parse()

	targetConfig, err := config.ReadTargets(*targetConfigPath)
	if err != nil {
		log.Fatalw("target config error", "path", *targetConfigPath, "error", err)
	}

	for _, target := range targetConfig.Targets {
		targetCollector, err := exporter.NewTarget(target)
		if err != nil {
			log.Fatalw("error creating target", "id", target.GatewayID, "baseUrl", target.BaseUrl, "apiKeyType", target.APIKey[:5])
		}
		err = prometheus.Register(targetCollector)
		if err != nil {
			log.Fatalw("error registering target", "id", target.GatewayID, "error", err)
		}
	}

	log.Infow("listening", "address", *address)
	srv := server.NewServer(*address)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalw("listening error", "addr", *address, "error", err)
	}
}
