package main

import (
	"github.com/bersennaidoo/arcbox/application/rest/handlers"
	"github.com/bersennaidoo/arcbox/application/rest/mid"
	"github.com/bersennaidoo/arcbox/application/rest/server"
	"github.com/bersennaidoo/arcbox/foundation/formdecode"
	"github.com/bersennaidoo/arcbox/infrastructure/repositories/mysql"
	"github.com/bersennaidoo/arcbox/physical/config"
	"github.com/bersennaidoo/arcbox/physical/dbc"
	"github.com/bersennaidoo/arcbox/physical/logger"
)

func main() {
	log := logger.New()
	filename := config.GetConfigFileName()
	cfg := config.New(filename)
	dbc := dbc.New(cfg)
	sr := mysql.New(dbc)
	decoder := formdecode.New()
	tempcache, _ := handlers.NewTemplateCache()
	h := handlers.New(log, sr, tempcache, decoder)
	m := mid.New(log)
	srv := server.New(h, cfg, log, m)
	srv.InitRouter()

	log.Info("Starting the application...")
	srv.Start()
}
