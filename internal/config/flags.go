package config

import (
	"github.com/spf13/pflag"
)

type AgentFlags struct {
	Address        string
	WriteInterval  int
	Collectnterval int
}

func InitAgentFlags() *AgentFlags {
	agentFlags := AgentFlags{}

	pflag.StringVarP(&agentFlags.Address, "addr", "a", "http://localhost:8080", "Address host:port")
	pflag.IntVarP(&agentFlags.WriteInterval, "reportInterval", "r", 10,
		"Wait interval in seconds before sending metrics to server")
	pflag.IntVarP(&agentFlags.Collectnterval, "pollInterval", "p", 2,
		"Wait interval in seconds before reading system metrics")

	pflag.Parse()

	return &agentFlags
}

type ServerFlags struct {
	Address        string
	ReportInterval int
	PollInterval   int
}

func InitServerFlags() *ServerFlags {
	serverFlags := ServerFlags{}

	pflag.StringVarP(&serverFlags.Address, "addr", "a", "localhost:8080", "Server address host:port")

	pflag.Parse()

	return &serverFlags
}
