package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	core_pgx_pool "github.com/Kriffi/golang-todoapp/internal/core/repository/postgres/pool/pgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func checkMigrationVersion(directory string) int {
	target, err := latestMigrationVersion(directory)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	config, err := core_pgx_pool.NewConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot configure the migration check:", err)
		return 2
	}
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	pool, err := core_pgx_pool.NewPool(ctx, config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot connect for the migration check:", err)
		return 2
	}
	defer pool.Close()

	var version uint64
	var dirty bool
	err = pool.Pool.QueryRow(ctx, "SELECT version, dirty FROM public.schema_migrations", pgx.QueryExecModeSimpleProtocol).Scan(&version, &dirty)
	var pgError *pgconn.PgError
	switch {
	case errors.Is(err, pgx.ErrNoRows), errors.As(err, &pgError) && pgError.Code == "42P01":
		return 1
	case err != nil:
		fmt.Fprintln(os.Stderr, "cannot read the database migration version:", err)
		return 2
	case dirty || version > target:
		fmt.Fprintln(os.Stderr, "database migration state does not match this build")
		return 2
	case version < target:
		return 1
	default:
		fmt.Printf("Database migrations are current (version %d)\n", version)
		return 0
	}
}

func latestMigrationVersion(directory string) (uint64, error) {
	files, err := filepath.Glob(filepath.Join(directory, "*.up.sql"))
	if err != nil {
		return 0, err
	}
	if len(files) == 0 {
		return 0, fmt.Errorf("no up migrations in %s", directory)
	}
	var latest uint64
	for _, file := range files {
		prefix, _, ok := strings.Cut(filepath.Base(file), "_")
		version, err := strconv.ParseUint(prefix, 10, 64)
		if !ok || err != nil {
			return 0, fmt.Errorf("invalid migration filename %s", filepath.Base(file))
		}
		latest = max(latest, version)
	}
	return latest, nil
}
