package model

type AuditLog struct {
	Ts        int64    `json:"ts"`
	Metrics   []string `json:"metrics"`
	IpAddress string   `json:"ip_address"`
}
