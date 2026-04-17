package postgres

import (
	"avanpost-show/pkg/config"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPgxConn pool
func NewPgxConn(cfg *config.Database) (*pgxpool.Pool, error) {
	ctx := context.Background()

	poolCfg, err := pgxpool.ParseConfig(cfg.ConnectionString)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = int32(cfg.PoolMaxConns)
	poolCfg.HealthCheckPeriod = time.Duration(cfg.PoolHealthCheckPeriodMilliseconds) * time.Millisecond
	poolCfg.MaxConnIdleTime = time.Duration(cfg.PoolMaxConnIdleTimeMilliseconds) * time.Millisecond
	poolCfg.MaxConnLifetime = time.Duration(cfg.PoolMaxConnLifetimeMilliseconds) * time.Millisecond
	poolCfg.MinConns = int32(cfg.PoolMinConns)

	connPool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("pgx.ConnectConfig %w", err)
	}

	return connPool, nil
}
