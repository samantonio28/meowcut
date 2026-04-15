package main

import (
	"flag"
	"log"

	"github.com/samantonio28/meowcut/internal/delivery"
)

func main() {
	configPath := flag.String("config", "configs/server.yaml", "path to configuration file")
	flag.Parse()

	app, err := delivery.NewApp(*configPath)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}