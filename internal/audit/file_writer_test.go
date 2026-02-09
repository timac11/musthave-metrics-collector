package audit

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/timac11/musthave-metrics-collector/internal/model"
)

func TestPositiveReadWriteWriterScenario(t *testing.T) {
	filePath := "/tmp/audit-writer.json"

	auditFileWriter := NewAuditLogFileWriter(filePath)

	nowTime1 := time.Now().Unix()

	auditLog1 := model.AuditLog{
		TS:        nowTime1,
		Metrics:   []string{"first"},
		IPAddress: "192.168.1.0",
	}

	err := auditFileWriter.write(&auditLog1)
	require.NoError(t, err)

	writtenLogs, err := auditFileWriter.read()

	require.NoError(t, err)
	assert.Equal(t, len(writtenLogs), 1)
	assert.Equal(t, writtenLogs[0].TS, auditLog1.TS)
	assert.Equal(t, writtenLogs[0].Metrics[0], auditLog1.Metrics[0])
	assert.Equal(t, writtenLogs[0].IPAddress, auditLog1.IPAddress)

	nowTime2 := time.Now().Unix()

	auditLog2 := model.AuditLog{
		TS:        nowTime2,
		Metrics:   []string{"second"},
		IPAddress: "192.168.1.1",
	}

	err = auditFileWriter.write(&auditLog2)
	writtenLogs, err = auditFileWriter.read()

	require.NoError(t, err)
	assert.Equal(t, 2, len(writtenLogs))
	assert.Equal(t, writtenLogs[0].TS, auditLog1.TS)
	assert.Equal(t, writtenLogs[0].Metrics[0], auditLog1.Metrics[0])
	assert.Equal(t, writtenLogs[0].IPAddress, auditLog1.IPAddress)

	assert.Equal(t, writtenLogs[1].TS, auditLog2.TS)
	assert.Equal(t, writtenLogs[1].Metrics[0], auditLog2.Metrics[0])
	assert.Equal(t, writtenLogs[1].IPAddress, auditLog2.IPAddress)

	defer os.Remove(filePath)
}
