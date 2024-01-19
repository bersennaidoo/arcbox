package main

import (
	"github.com/bersennaidoo/arcbox/foundation/formdecode"
	"github.com/bersennaidoo/arcbox/infrastructure/repositories/mysql"
	"github.com/bersennaidoo/arcbox/physical/config"
	"github.com/bersennaidoo/arcbox/physical/dbc"
	"github.com/bersennaidoo/arcbox/physical/logger"
	"github.com/bersennaidoo/arcbox/physical/session"
	"github.com/bersennaidoo/arcbox/transport/http/application"
	"github.com/bersennaidoo/arcbox/transport/http/handlers"
	"github.com/bersennaidoo/arcbox/transport/http/mid"
)

func main() {
	log := logger.New()
	filename := config.GetConfigFileName()
	cfg := config.New(filename)
	dbc := dbc.New(cfg)
	sr := mysql.NewSnipsRepository(dbc)
	ur := mysql.NewUsersRepository(dbc)
	decoder := formdecode.New()
	tempcache, _ := handlers.NewTemplateCache()
	sessionM := session.New(dbc)
	h := handlers.New(log, sr, ur, tempcache, decoder, sessionM)
	m := mid.New(log, sessionM, ur)
	app := application.New(h, cfg, log, m, sessionM)
	n := app.InitRouter()

	log.Info("Starting the application...")
	app.Start(n)
}
