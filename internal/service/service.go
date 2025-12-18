package service

import (
	"context"
	"errors"
	"time"

	"github.com/avast/retry-go/v4"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/timac11/musthave-metrics-collector/internal/model"
)

type Repository interface {
	Save(ctx context.Context, value model.Metrics) error
	Get(ctx context.Context, id string, mType string) (*model.Metrics, error)
	GetAll(ctx context.Context) ([]model.Metrics, error)
	SaveAll(ctx context.Context, metrics []model.Metrics) error
	Ping(ctx context.Context) error
}

type Service struct {
	storage Repository
}

func NewService(storage Repository) *Service {
	service := &Service{
		storage: storage,
	}

	return service
}

func (service *Service) DBPing() error {
	return service.storage.Ping(context.Background())
}

func (service *Service) Save(metric model.Metrics) error {
	return retry.Do(
		func() error {
			storage := service.storage
			return storage.Save(context.Background(), metric)
		},
		service.getRetryOptions()...,
	)
}

func (service *Service) SaveAll(metrics []model.Metrics) error {
	return retry.Do(
		func() error {
			storage := service.storage
			return storage.SaveAll(context.Background(), metrics)
		},
		service.getRetryOptions()...,
	)
}

func (service *Service) Get(id string, mType string) (*model.Metrics, error) {
	var metric *model.Metrics
	var err error

	err = retry.Do(
		func() error {
			storage := service.storage
			metric, err = storage.Get(context.Background(), id, mType)
			return err
		},
		service.getRetryOptions()...,
	)

	return metric, err
}

func (service *Service) GetAll() ([]model.Metrics, error) {
	var metrics []model.Metrics
	var err error

	err = retry.Do(
		func() error {
			storage := service.storage
			metrics, err = storage.GetAll(context.Background())
			return err
		},
		service.getRetryOptions()...,
	)

	return metrics, err
}

func (service *Service) getRetryOptions() []retry.Option {
	return []retry.Option{
		retry.RetryIf(func(errAttempt error) bool {
			var pgErr *pgconn.PgError
			if errors.As(errAttempt, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
				return true
			}
			return false
		}),
		retry.Attempts(3),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			return time.Second + time.Duration(n*2)*time.Second
		}),
		retry.Context(context.Background()),
	}
}
