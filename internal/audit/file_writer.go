package audit

import (
	"encoding/json"
	"os"

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

		logs, err := writer.read()
		if err != nil {
			logs = []*model.AuditLog{}
		}

		logs = append(logs, observable.value)

		err = writer.write(logs)

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

func (writer *AuditLogFileWriter) write(logs []*model.AuditLog) error {
	file, err := os.OpenFile(writer.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	defer file.Close()

	data, err := json.Marshal(&logs)

	if err != nil {
		return err
	}

	file.Write(data)

	return nil
}
