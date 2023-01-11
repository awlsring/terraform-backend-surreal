package main

import (
	"log"

	"github.com/awlsring/surreal-db-client/surreal"
	"github.com/awlsring/terraform-backend-surreal/pkg/config"
	"github.com/awlsring/terraform-backend-surreal/pkg/server"
	"github.com/awlsring/terraform-backend-surreal/pkg/state"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	s, err := surreal.New(cfg.Surreal)
	if err != nil {
		log.Fatalln(err)
	}

	dao := state.New(s)

	server.Start(&cfg, dao)
}