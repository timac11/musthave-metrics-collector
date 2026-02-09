package audit

import (
	"encoding/json"
	"os"

	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

type AuditLogFileWriter struct {
	path string
}

func NewAuditLogFileWriter(path string) *AuditLogFileWriter {
	return &AuditLogFileWriter{path: path}
}

func (writer *AuditLogFileWriter) Subscribe(observable *Observable) {
	for {
		cond := observable.cond
		cond.L.Lock()
		cond.Wait()

		err := writer.write(observable.value)

		if err != nil {
			logger.Error("Failed to write metrics", err)
		}

		cond.L.Unlock()
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
		logger.Error("Failed to read metrics", err)
		existedLogs = []*model.AuditLog{}
	}

	existedLogs = append(existedLogs, log)

	data, err := json.Marshal(existedLogs)
	if err != nil {
		return err
	}

	logger.Error("Info metric", len(existedLogs))
	logger.Error("Info metric", string(data))

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
