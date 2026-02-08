package audit

import (
	"strings"

	"github.com/go-resty/resty/v2"
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

		serverSender.client.R().SetBody(observable.value).Post("/")

		cond.L.Unlock()
	}
}
