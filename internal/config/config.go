package config

import (
	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
)

type AgentConfig struct {
	Address        string `env:"ADDRESS"`
	WriteInterval  int    `env:"REPORT_INTERVAL"`
	Collectnterval int    `env:"POLL_INTERVAL"`
}

func InitAgentConfig() *AgentConfig {
	envConfig := initAgentEnv()
	flagsConfig := initAgentFlags()

	if envConfig.Address == "" {
		envConfig.Address = flagsConfig.Address
	}

	if envConfig.Collectnterval == 0 {
		envConfig.Collectnterval = flagsConfig.Collectnterval
	}

	if envConfig.WriteInterval == 0 {
		envConfig.WriteInterval = flagsConfig.WriteInterval
	}

	return envConfig
}

func initAgentFlags() *AgentConfig {
	agentConfig := AgentConfig{}

	pflag.StringVarP(&agentConfig.Address, "addr", "a", "http://localhost:8080", "Address host:port")
	pflag.IntVarP(&agentConfig.WriteInterval, "reportInterval", "r", 10,
		"Wait interval in seconds before sending metrics to server")
	pflag.IntVarP(&agentConfig.Collectnterval, "pollInterval", "p", 2,
		"Wait interval in seconds before reading system metrics")

	pflag.Parse()

	return &agentConfig
}

func initAgentEnv() *AgentConfig {
	agentConfig := AgentConfig{}

	env.Parse(&agentConfig)

	return &agentConfig
}

type ServerConfig struct {
	Address string `env:"ADDRESS"`
}

func InitServerConfig() *ServerConfig {
	serverEnv := initServerEnv()

	if serverEnv.Address != "" {
		return serverEnv
	}

	return initServerFlags()
}

func initServerEnv() *ServerConfig {
	serverConfig := ServerConfig{}

	env.Parse(&serverConfig)

	return &serverConfig
}

func initServerFlags() *ServerConfig {
	serverFlags := ServerConfig{}

	pflag.StringVarP(&serverFlags.Address, "addr", "a", "localhost:8080", "Server address host:port")

	pflag.Parse()

	return &serverFlags
}
