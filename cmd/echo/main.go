package main

import (
	"log/slog"

	"github.com/alexflint/go-arg"
	"github.com/k0kubun/pp/v3"
	"github.com/zekroTJA/echo/pkg/server"
	"github.com/zekroTJA/echo/pkg/verbosity"
	"github.com/zekrotja/parsables"
	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/level"
)

type Config struct {
	Address   string              `arg:"-a,--address,env:ECHO_ADDRESS" default:":80" help:"Bind address"`
	Verbosity verbosity.Verbosity `arg:"-v,--verbosity,env:ECHO_VERBOSITY" default:"normal" help:"Request response verbosity"`
	BodyLimit parsables.FileSize  `arg:"-b,--body-limit,env:ECHO_BODY_LIMIT" default:"512kib" help:"Maximum bytes to read from request body"`
	LogLevel  level.Level         `arg:"-l,--log-level,env:ECHO_LOG_LEVEL" default:"info" help:"Log level"`
}

func main() {
	var cfg Config
	arg.MustParse(&cfg)

	logger := slog.New(rogu.NewLogger(rogu.NewPrettyWriter()).SetLevel(cfg.LogLevel))
	slog.SetDefault(logger)

	slog.Debug("loaded config", "config", pp.Sprint(cfg))

	s := server.New(cfg.Address, cfg.Verbosity, int(cfg.BodyLimit))
	slog.Info("starting server", "address", cfg.Address)
	if err := s.Run(); err != nil {
		slog.Error("failed starting server", "err", err)
	}
}
