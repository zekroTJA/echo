package main

import (
	"log"

	"github.com/alexflint/go-arg"
	"github.com/k0kubun/pp/v3"
	"github.com/zekroTJA/echo/pkg/server"
	"github.com/zekroTJA/echo/pkg/verbosity"
)

type Config struct {
	Address   string              `arg:"-a,--address,env:ECHO_ADDRESS" default:":80"`
	Verbosity verbosity.Verbosity `arg:"-v,--verbosity,env:ECHO_VERBOSITY" default:"normal"`
}

func main() {
	var cfg Config
	arg.MustParse(&cfg)

	pp.Println(cfg)

	s := server.New(cfg.Address, cfg.Verbosity)
	log.Printf("Running server on address %s...", cfg.Address)
	if err := s.Run(); err != nil {
		log.Fatalf("Failed starting server: %s", err.Error())
	}
}
