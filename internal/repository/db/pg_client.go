package dbstorage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
	"github.com/timac11/musthave-metrics-collector/internal/model"
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

func (client *PgClient) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if client.conn == nil {
		return errors.New("connection was not established")
	}

	return client.conn.Ping(ctx)
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

func (client * PgClient) Get(ctx context.Context, id string, mType string) (*model.Metrics, error) {
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
    
    err = client.conn.QueryRow(ctx, query, id, mType).Scan(
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
    
    rows, err := client.conn.Query(ctx, query)
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
