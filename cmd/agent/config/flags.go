package config

import (
	"github.com/spf13/pflag"
)

type AgentFlags struct {
	Address        string
	ReportInterval int
	PollInterval   int
}

func InitFlags() *AgentFlags {
	agentFlags := AgentFlags{}
	pflag.StringVarP(&agentFlags.Address, "addr", "a", "http://localhost:8080", "Address host:port")
	pflag.IntVarP(&agentFlags.ReportInterval, "reportInterval", "r", 10,
		"Wait interval in seconds before sending metrics to server")
	pflag.IntVarP(&agentFlags.PollInterval, "pollInterval", "p", 2,
		"Wait interval in seconds before reading system metrics")

	pflag.Parse()

	return &agentFlags
}
