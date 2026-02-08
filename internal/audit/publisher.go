package audit

import (
	"sync"

	"github.com/timac11/musthave-metrics-collector/internal/model"
)

type Observable struct {
	cond  *sync.Cond
	value *model.AuditLog
}

type AuditLogsPublisher struct {
	observable *Observable
}

type Subscriber interface {
	Subscribe(observable *Observable)
}

func NewAuditLogsPublisher() *AuditLogsPublisher {
	observable := &Observable{
		cond: sync.NewCond(new(sync.Mutex)),
	}

	return &AuditLogsPublisher{
		observable: observable,
	}
}

func (publisher *AuditLogsPublisher) Register(subscriber Subscriber) {
	go subscriber.Subscribe(publisher.observable)
}

func (publisher *AuditLogsPublisher) Publish(log *model.AuditLog) {
	observable := publisher.observable
	observable.cond.L.Lock()
	observable.value = log

	publisher.observable.cond.Broadcast()
	observable.cond.L.Unlock()
}
