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
	CryptoKey      string `env:"CRYPTO_KEY"`
	Config         string `env:"CONFIG"`
	RetryAttempts  uint
	RetryInterval  uint
}

type agentJsonConfig struct {
	Address        string `json:"address"`
	ReportInterval int    `json:"report_interval"`
	PollInterval   int    `json:"poll_interval"`
	CryptoKey      string `json:"crypto_key"`
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

	if agentEnv.CryptoKey == "" {
		agentEnv.CryptoKey = agentFlags.CryptoKey
	}

	if agentEnv.Config == "" {
		agentEnv.Config = agentFlags.Config
	}

	if agentEnv.Config != "" {
		// TODO: add variables from json file
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
	pflag.StringVar(&agentFlags.CryptoKey, "crypto-key", "", "Public key file path")
	pflag.StringVarP(&agentFlags.Config, "config", "c", "", "Path to config json")

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
	AuditURL        string `env:"AUDIT_URL"`
	CryptoKey       string `env:"CRYPTO_KEY"`
	Config          string `env:"CONFIG"`
	RetryAttempts   uint
	RetryInterval   uint
}

type serverJsonConfig struct {
	Address         string `json:"address"`
	Restore         bool   `json:"restore"`
	FileStoragePath string `json:"store_file"`
	DatabaseDsn     string `json:"database_dsn"`
	CryptoKey       string `json:"crypto_key"`
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

	if serverEnv.AuditURL == "" {
		serverEnv.AuditURL = serverFlags.AuditURL
	}

	if serverEnv.SigningKey == "" {
		serverEnv.SigningKey = serverFlags.SigningKey
	}

	if !serverEnv.Restore {
		serverEnv.Restore = serverFlags.Restore
	}

	if serverEnv.CryptoKey == "" {
		serverEnv.CryptoKey = serverFlags.CryptoKey
	}

	if serverEnv.Config == "" {
		serverEnv.Config = serverFlags.Config
	}

	serverEnv.RetryAttempts = serverFlags.RetryAttempts
	serverEnv.RetryInterval = serverFlags.RetryInterval

	if serverEnv.Config != "" {
		// TODO: assign variables from file
	}

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
	pflag.StringVar(&serverFlags.AuditURL, "audit-url", "", "Url to send audit logs")
	pflag.StringVar(&serverFlags.CryptoKey, "crypto-key", "", "Private key file path")
	pflag.StringVarP(&serverFlags.Config, "config", "c", "", "Path to config json")

	pflag.Parse()

	return &serverFlags
}
