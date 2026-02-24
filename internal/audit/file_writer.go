package audit

import (
	"encoding/json"
	"os"

	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

// AuditLogFileWriter is structure is used to write logs to file
type AuditLogFileWriter struct {
	path string
}

// NewAuditLogFileWriter - return instance of AuditLogFileWriter.
// implements interface Subscriber
func NewAuditLogFileWriter(path string) *AuditLogFileWriter {
	return &AuditLogFileWriter{path: path}
}

// Subscribe is called for subscribing on event of writing new audit logs
func (writer *AuditLogFileWriter) Subscribe(observable *Observable) {
	logger.Info("subscribe to write logs to file")

	for {
		observable.cond.L.Lock()
		observable.cond.Wait()

		logger.Info("start write audit logs to file")
		err := writer.write(observable.value)

		if err != nil {
			logger.Error("failed to write metrics", err)
		}

		observable.cond.L.Unlock()
	}
}

func (writer *AuditLogFileWriter) read() ([]*model.AuditLog, error) {
	data, err := os.ReadFile(writer.path)
	if err != nil {
		return nil, err
	}

	var logs []*model.AuditLog
	err = json.Unmarshal(data, &logs)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (writer *AuditLogFileWriter) write(log *model.AuditLog) error {
	file, err := os.OpenFile(writer.path, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer file.Close()

	existedLogs, err := writer.read()
	if err != nil {
		existedLogs = []*model.AuditLog{}
	}

	existedLogs = append(existedLogs, log)
	data, err := json.Marshal(existedLogs)

	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
