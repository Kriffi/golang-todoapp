package core_pgx_pool

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	core_postgres_pool "github.com/Kriffi/golang-todoapp/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func NewPool(
	ctx context.Context,
	config Config,
) (*Pool, error) {
	connectionString := config.DatabaseURL
	if connectionString == "" {
		localURL := url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword(config.User, config.Password),
			Host:     net.JoinHostPort(config.Host, config.Port),
			Path:     "/" + config.Database,
			RawQuery: "sslmode=disable",
		}
		connectionString = localURL.String()
	}

	pgxconfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("invalid PostgreSQL connection configuration")
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pgxpool ping: %w", err)
	}

	return &Pool{
		Pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

func (p *Pool) Query(
	ctx context.Context,
	sql string,
	args ...any,
) (core_postgres_pool.Rows, error) {
	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, mapErrors(err)
	}

	return pgxRows{Rows: rows}, nil
}

func (p *Pool) QueryRow(
	ctx context.Context,
	sql string,
	args ...any,
) core_postgres_pool.Row {
	return pgxRow{Row: p.Pool.QueryRow(ctx, sql, args...)}
}

func (p *Pool) Exec(
	ctx context.Context,
	sql string,
	arguments ...any,
) (core_postgres_pool.CommandTag, error) {
	cmdTag, err := p.Pool.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, mapErrors(err)
	}

	return pgxCommandTag{CommandTag: cmdTag}, nil
}

func (p *Pool) OpTimeout() time.Duration {
	return p.opTimeout
}
