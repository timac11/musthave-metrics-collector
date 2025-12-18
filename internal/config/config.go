package config

import (
	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
)

type AgentConfig struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	RetryAttempts  uint
	RetryInterval  uint
}

func InitAgentConfig() *AgentConfig {
	envConfig := initAgentEnv()
	flagsConfig := initAgentFlags()

	if envConfig.Address == "" {
		envConfig.Address = flagsConfig.Address
	}

	if envConfig.PollInterval == 0 {
		envConfig.PollInterval = flagsConfig.PollInterval
	}

	if envConfig.ReportInterval == 0 {
		envConfig.ReportInterval = flagsConfig.ReportInterval
	}

	return envConfig
}

func initAgentFlags() *AgentConfig {
	agentFlags := AgentConfig{}

	pflag.StringVarP(&agentFlags.Address, "addr", "a", "http://localhost:8080", "Address host:port")
	pflag.IntVarP(&agentFlags.ReportInterval, "reportInterval", "r", 10, "Wait interval in seconds before sending metrics to server")
	pflag.IntVarP(&agentFlags.PollInterval, "pollInterval", "p", 2, "Wait interval in seconds before reading system metrics")
	pflag.UintVar(&agentFlags.RetryAttempts, "retryAttempt", 3, "Count of retry attempts to execute metrics operation")
	pflag.UintVar(&agentFlags.RetryInterval, "retryInterval", 2, "Interval in seconds between metric operation attempts")

	pflag.Parse()

	return &agentFlags
}

func initAgentEnv() *AgentConfig {
	agentConfig := AgentConfig{}

	env.Parse(&agentConfig)

	return &agentConfig
}

type ServerConfig struct {
	Address         string `env:"ADDRESS"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDsn     string `env:"DATABASE_DSN"`
	RetryAttempts   uint
	RetryInterval   uint
}

func InitServerConfig() *ServerConfig {
	serverEnv := initServerEnv()
	serverFlags := initServerFlags()

	if serverEnv.Address == "" {
		serverEnv.Address = serverFlags.Address
	}

	if serverEnv.FileStoragePath == "" {
		serverEnv.FileStoragePath = serverFlags.FileStoragePath
	}

	if serverEnv.DatabaseDsn == "" {
		serverEnv.DatabaseDsn = serverFlags.DatabaseDsn
	}

	if !serverEnv.Restore {
		serverEnv.Restore = serverFlags.Restore
	}

	serverEnv.RetryAttempts = serverFlags.RetryAttempts
	serverEnv.RetryInterval = serverFlags.RetryInterval

	return serverEnv
}

func initServerEnv() *ServerConfig {
	serverConfig := ServerConfig{}

	env.Parse(&serverConfig)

	return &serverConfig
}

func initServerFlags() *ServerConfig {
	serverFlags := ServerConfig{}

	pflag.StringVarP(&serverFlags.Address, "addr", "a", "localhost:8080", "Server address host:port")
	pflag.StringVarP(&serverFlags.DatabaseDsn, "db", "d", "", "Postgres database url")
	pflag.StringVarP(&serverFlags.FileStoragePath, "file", "f", "./db.json", "File to store JSON file with metrics")
	pflag.BoolVarP(&serverFlags.Restore, "restore", "r", true, "Restore or not metrics from file storage")
	pflag.UintVar(&serverFlags.RetryAttempts, "retryAttempt", 3, "Count of retry attempts to execute metrics operation")
	pflag.UintVar(&serverFlags.RetryInterval, "retryInterval", 2, "Interval in seconds between metric operation attempts")

	pflag.Parse()

	return &serverFlags
}
