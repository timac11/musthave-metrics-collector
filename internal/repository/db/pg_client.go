package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
)

type PgClient struct {
	conn *sql.DB
}

func NewPgClient(url string) *PgClient {
	conn, err := sql.Open("pgx", url)

	if err != nil {
		logger.Error("Failed connect to database")
		logger.Error(err.Error())
		log.Fatal(err)
	}

	client := &PgClient{conn: conn}

	err = client.applyMigration(context.Background())

	return client
}

func (client *PgClient) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if client.conn == nil {
		return errors.New("connection was not established")
	}

	return client.conn.Ping()
}

func (client *PgClient) Save(ctx context.Context, metric model.Metrics) error {
	err := client.Ping(ctx)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO metrics (name, mtype, delta, value, hash)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (name, mtype) 
        DO UPDATE SET
            delta = EXCLUDED.delta,
            value = EXCLUDED.value,
            hash = EXCLUDED.hash,
            updated_at = CURRENT_TIMESTAMP
    `

	_, err = client.conn.Exec(
		query,
		metric.ID,
		metric.MType,
		metric.Delta,
		metric.Value,
		metric.Hash,
	)

	return err
}

func (client *PgClient) Get(ctx context.Context, id string, mType string) (*model.Metrics, error) {
	err := client.Ping(ctx)
	if err != nil {
		return nil, err
	}

	query := `
        SELECT name, mtype, delta, value, hash
        FROM metrics
        WHERE name = $1 AND mtype = $2
    `

	var metric model.Metrics

	err = client.conn.QueryRow(query, id, mType).Scan(
		&metric.ID,
		&metric.MType,
		&metric.Delta,
		&metric.Value,
		&metric.Hash,
	)

	if err != nil {
		logger.Error("Failed to get metric", err)
		return nil, err
	}

	return &metric, nil
}

func (client *PgClient) GetAll(ctx context.Context) ([]model.Metrics, error) {
	err := client.Ping(ctx)
	if err != nil {
		return nil, err
	}

	query := `
        SELECT name, mtype, delta, value, hash
        FROM metrics
        ORDER BY name, mtype
    `

	rows, err := client.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []model.Metrics
	for rows.Next() {
		var metric model.Metrics

		err := rows.Scan(
			&metric.ID,
			&metric.MType,
			&metric.Delta,
			&metric.Value,
			&metric.Hash,
		)

		if err != nil {
			return nil, err
		}

		metrics = append(metrics, metric)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return metrics, nil
}

func (client *PgClient) applyMigration(ctx context.Context) error {
	err := client.Ping(ctx)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(client.conn, &postgres.Config{})
	if err != nil {
		logger.Error("failed to create driver", err)
		return err
	}

	migrations, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres",
		driver,
	)

	if err != nil {
		logger.Error("failed to apply migrations", err)
		return err
	}

	// Apply migrations
	err = migrations.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Error("failed to apply migrations", err)
		return err
	}

	logger.Info("migrations successfully applied")

	return nil
}
