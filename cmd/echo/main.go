package main

import (
	"log"

	"github.com/spf13/viper"
	"github.com/zekroTJA/echo/pkg/config"
	"github.com/zekroTJA/echo/pkg/server"
	"github.com/zekroTJA/echo/pkg/verbosity"
)

func main() {
	config.InitViper()

	addr := viper.GetString(config.KeyAddr)
	debug := viper.GetBool(config.KeyDebug)
	verb, err := verbosity.FromInt(viper.GetInt(config.KeyVerbosity))
	if err != nil {
		log.Fatalf("Startup failed: %s", err.Error())
	}

	s := server.New(addr, verb, debug)
	log.Printf("Running server on address %s...", addr)
	if err = s.Run(); err != nil {
		log.Fatalf("Failed starting server: %s", err.Error())
	}
}
