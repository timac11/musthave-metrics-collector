package audit

import (
	"github.com/timac11/musthave-metrics-collector/internal/model"
	"time"
)

type Auditor struct {
	publisher *AuditLogsPublisher
}

func NewAuditor(auditsPath string, auditsUrl string) *Auditor {
	publisher := NewAuditLogsPublisher()

	if auditsPath != "" {
		fileWriter := NewAuditLogFileWriter(auditsPath)
		publisher.Register(fileWriter)
	}

	if auditsUrl != "" {
		serverSender := NewAuditLogServerSender(auditsUrl)
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

	log := model.AuditLog{Ts: now, Metrics: metricNames, IpAddress: ipAddress}

	go auditor.publisher.Publish(&log)
}
