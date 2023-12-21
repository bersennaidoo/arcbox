package main

import (
	"github.com/bersennaidoo/arcbox/application/rest/handlers"
	"github.com/bersennaidoo/arcbox/application/rest/server"
)

func main() {
	h := handlers.New()
	srv := server.New(h)
	srv.InitRouter()
	srv.Start()
}
