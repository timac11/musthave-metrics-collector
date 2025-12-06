package dbstorage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
)

type PgClient struct {
	conn *pgx.Conn
}

func NewPgClient(url string) *PgClient {
	conn, err := pgx.Connect(context.Background(), url)

	if err != nil {
		logger.Error("Failed connect to database")
		logger.Error(err.Error())
	}

	client := &PgClient{conn: conn}
	return client
}

func (client *PgClient) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if client.conn == nil {
		return errors.New("connection was not estableshed")
	}

	return client.conn.Ping(ctx)
}
