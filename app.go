package main

import (
	"github.com/pikomonde/pikopos/clients"
	"github.com/pikomonde/pikopos/delivery"
	"github.com/pikomonde/pikopos/repository"
	"github.com/pikomonde/pikopos/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	// setup client
	cli := clients.New()

	// setup repository
	repo := repository.New(cli)

	// setup service
	serv := service.New(repo)

	// setup delivery
	dlvr := delivery.New(serv)
	dlvr.Start()
}
