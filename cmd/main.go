package main

import (
	"github.com/awlsring/surreal-db-client/surreal"
	"github.com/awlsring/terraform-backend-surreal/pkg/config"
	"github.com/awlsring/terraform-backend-surreal/pkg/server"
	"github.com/awlsring/terraform-backend-surreal/pkg/state"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting Terraform Surreal Backend")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error(err)
	}

	s, err := surreal.New(cfg.Surreal)
	if err != nil {
		log.Error(err)
	}

	dao := state.New(s)

	server.Start(&cfg, dao)
}
