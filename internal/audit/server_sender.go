package audit

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

type AuditLogServerSender struct {
	client resty.Client
}

func NewAuditLogServerSender(url string) *AuditLogServerSender {
	client := resty.New()

	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	client.SetBaseURL(url)

	return &AuditLogServerSender{client: *client}
}

func (serverSender *AuditLogServerSender) Subscribe(observable *Observable) {
	for {
		cond := observable.cond
		cond.L.Lock()
		cond.Wait()

		err := serverSender.send(observable.value)

		if err != nil {
			logger.Error("Failed to write metrics", err)
		}

		cond.L.Unlock()
	}
}

func (serverSender *AuditLogServerSender) send(log *model.AuditLog) error {
	response, err := serverSender.client.R().SetBody(log).Post("")

	if err != nil {
		return err
	}

	if !response.IsSuccess() {
		return fmt.Errorf("Failed send metrics audit log, status code: %d", response.StatusCode())
	}

	return nil
}
