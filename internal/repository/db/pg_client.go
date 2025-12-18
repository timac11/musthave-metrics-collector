package dbstorage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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

func NewPgClient(url string) (*PgClient, error) {
	conn, err := sql.Open("pgx", url)

	if err != nil {
		logger.Error("Failed connect to database", err.Error())
		return nil, err
	}

	client := &PgClient{conn: conn}

	err = client.applyMigration(context.Background())

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (client *PgClient) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	return client.conn.PingContext(ctx)
}

func (client *PgClient) Save(ctx context.Context, metric model.Metrics) error {
	query := `
    INSERT INTO metrics (name, mtype, delta, value, hash)
    VALUES ($1, $2, $3, $4, $5)
    ON CONFLICT (name, mtype) 
    DO UPDATE SET
        delta = CASE 
            WHEN EXCLUDED.mtype = 'counter' AND EXCLUDED.delta IS NOT NULL
            THEN COALESCE(metrics.delta, 0) + EXCLUDED.delta
            ELSE EXCLUDED.delta
        END,
        value = EXCLUDED.value,
        hash = EXCLUDED.hash,
        updated_at = CURRENT_TIMESTAMP
    `

	_, err := client.conn.ExecContext(
		ctx,
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
	query := `
        SELECT name, mtype, delta, value, hash
        FROM metrics
        WHERE name = $1 AND mtype = $2
    `

	var metric model.Metrics

	err := client.conn.QueryRowContext(ctx, query, id, mType).Scan(
		&metric.ID,
		&metric.MType,
		&metric.Delta,
		&metric.Value,
		&metric.Hash,
	)

	if err != nil {
		return nil, err
	}

	return &metric, nil
}

func (client *PgClient) SaveAll(ctx context.Context, metrics []model.Metrics) error {
	if len(metrics) == 0 {
		return nil
	}

	if err := client.Ping(ctx); err != nil {
		return err
	}

	tx, err := client.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	valueStrings := make([]string, 0, len(metrics))
	valueArgs := make([]interface{}, 0, len(metrics)*5)
	index := 1

	for _, metric := range metrics {

		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
				index, index+1, index+2, index+3, index+4))

		valueArgs = append(valueArgs, metric.ID)
		valueArgs = append(valueArgs, metric.MType)
		valueArgs = append(valueArgs, metric.Delta)
		valueArgs = append(valueArgs, metric.Value)
		valueArgs = append(valueArgs, metric.Hash)

		index += 5
	}

	query := fmt.Sprintf(`
    	INSERT INTO metrics (name, mtype, delta, value, hash)
    	VALUES %s
    	ON CONFLICT (name, mtype) 
    	DO UPDATE SET
        	delta = CASE 
            	WHEN EXCLUDED.mtype = 'counter' AND EXCLUDED.delta IS NOT NULL
            	THEN COALESCE(metrics.delta, 0) + EXCLUDED.delta
            	ELSE EXCLUDED.delta
        	END,
        	value = EXCLUDED.value,
        	hash = EXCLUDED.hash,
        	updated_at = CURRENT_TIMESTAMP
    `, strings.Join(valueStrings, ","))

	_, err = tx.ExecContext(ctx, query, valueArgs...)

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
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

	rows, err := client.conn.QueryContext(ctx, query)
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
		return err
	}

	migrations, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres",
		driver,
	)

	if err != nil {
		return err
	}

	// Apply migrations
	err = migrations.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	logger.Info("migrations successfully applied")

	return nil
}
