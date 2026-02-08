package config

import (
	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
)

type AgentConfig struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	SigningKey     string `env:"KEY"`
	RateLimit      uint   `env:"RATE_LIMIT"`
	RetryAttempts  uint
	RetryInterval  uint
}

func InitAgentConfig() *AgentConfig {
	agentEnv := initAgentEnv()
	agentFlags := initAgentFlags()

	if agentEnv.Address == "" {
		agentEnv.Address = agentFlags.Address
	}

	if agentEnv.PollInterval == 0 {
		agentEnv.PollInterval = agentFlags.PollInterval
	}

	if agentEnv.ReportInterval == 0 {
		agentEnv.ReportInterval = agentFlags.ReportInterval
	}

	if agentEnv.RateLimit == 0 {
		agentEnv.RateLimit = agentFlags.RateLimit
	}

	if agentEnv.SigningKey == "" {
		agentEnv.SigningKey = agentFlags.SigningKey
	}

	return agentEnv
}

func initAgentFlags() *AgentConfig {
	agentFlags := AgentConfig{}

	pflag.StringVarP(&agentFlags.Address, "addr", "a", "http://localhost:8080", "Address host:port")
	pflag.IntVarP(&agentFlags.ReportInterval, "reportInterval", "r", 10, "Wait interval in seconds before sending metrics to server")
	pflag.IntVarP(&agentFlags.PollInterval, "pollInterval", "p", 2, "Wait interval in seconds before reading system metrics")
	pflag.UintVar(&agentFlags.RetryAttempts, "retryAttempt", 3, "Count of retry attempts to execute metrics operation")
	pflag.UintVar(&agentFlags.RetryInterval, "retryInterval", 2, "Interval in seconds between metric operation attempts")
	pflag.UintVarP(&agentFlags.RateLimit, "rateLimit", "l", 1, "Count of workers")

	pflag.StringVarP(&agentFlags.SigningKey, "signingKey", "k", "", "Signing key")

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
	SigningKey      string `env:"KEY"`
	AuditFile       string `env:"AUDIT_FILE"`
	AuditUrl        string `env:"AUDIT_URL"`
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

	if serverEnv.AuditFile == "" {
		serverEnv.AuditFile = serverFlags.AuditFile
	}

	if serverEnv.AuditUrl == "" {
		serverEnv.AuditUrl = serverFlags.AuditUrl
	}

	if serverEnv.SigningKey == "" {
		serverEnv.SigningKey = serverFlags.SigningKey
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
	pflag.UintVar(&serverFlags.RetryAttempts, "retry-attempt", 3, "Count of retry attempts to execute metrics operation")
	pflag.UintVar(&serverFlags.RetryInterval, "retry-interval", 2, "Interval in seconds between metric operation attempts")
	pflag.StringVarP(&serverFlags.SigningKey, "signing-key", "k", "", "Signing key")
	pflag.StringVar(&serverFlags.AuditFile, "audit-file", "", "File to store audit logs")
	pflag.StringVar(&serverFlags.AuditUrl, "audit-url", "", "Url to send audit logs")

	pflag.Parse()

	return &serverFlags
}
