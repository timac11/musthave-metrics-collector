package config

import (
	"github.com/spf13/pflag"
)

type ServerFlags struct {
	Address        string
	ReportInterval int
	PollInterval   int
}

func InitFlags() *ServerFlags {
	serverFlags := ServerFlags{}
	pflag.StringVarP(&serverFlags.Address, "addr", "a", "localhost:8080", "Server address host:port")

	pflag.Parse()

	return &serverFlags
}
