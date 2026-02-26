package audit

import (
	"time"

	"github.com/timac11/musthave-metrics-collector/internal/model"
)

// Auditor is the base structure to collect audit logs
// write audit logs to file and send logs to server if according parameters are set
type Auditor struct {
	publisher *AuditLogsPublisher // publisher is used to publish logs for subscribers
}

// NewAuditor creates new auditor instance
// auditsPath - is file path tp json file with audit logs
// auditsURL - server URL, where audit logs are send
func NewAuditor(auditsPath, auditsURL string) *Auditor {
	publisher := NewAuditLogsPublisher()

	if auditsPath != "" {
		fileWriter := NewAuditLogFileWriter(auditsPath)
		publisher.Register(fileWriter)
	}

	if auditsURL != "" {
		serverSender := NewAuditLogServerSender(auditsURL)
		publisher.Register(serverSender)
	}

	return &Auditor{publisher: publisher}
}

// Collect function called for writing audit logs
// metrics - metrics slice, from that audit logs are calculated
// ipAddress - is ip address of agent, which sent metrics
func (auditor *Auditor) Collect(metrics []*model.Metrics, ipAddress string) {
	now := time.Now().Unix()
	metricNames := make([]string, len(metrics))
	for i := range metrics {
		metricNames[i] = metrics[i].ID
	}

	log := model.AuditLog{TS: now, Metrics: metricNames, IPAddress: ipAddress}

	auditor.publisher.Publish(&log)
}
