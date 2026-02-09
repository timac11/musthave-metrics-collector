package audit

import (
	"time"

	"github.com/timac11/musthave-metrics-collector/internal/model"
)

type Auditor struct {
	publisher *AuditLogsPublisher
}

func NewAuditor(auditsPath string, auditsURL string) *Auditor {
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

func (auditor *Auditor) Collect(metrics []*model.Metrics, ipAddress string) {
	now := time.Now().Unix()
	metricNames := make([]string, len(metrics))
	for i := range metrics {
		metricNames[i] = metrics[i].MType
	}

	log := model.AuditLog{TS: now, Metrics: metricNames, IPAddress: ipAddress}

	auditor.publisher.Publish(&log)
}
