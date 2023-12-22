package main

import (
	"github.com/bersennaidoo/arcbox/application/rest/handlers"
	"github.com/bersennaidoo/arcbox/application/rest/server"
	"github.com/bersennaidoo/arcbox/physical/config"
	"github.com/bersennaidoo/arcbox/physical/logger"
)

func main() {
	log := logger.New()
	filename := config.GetConfigFileName()
	cfg := config.New(filename)
	h := handlers.New(log)
	srv := server.New(h, cfg, log)
	srv.InitRouter()

	log.Info("Starting the application...")
	srv.Start()
}
