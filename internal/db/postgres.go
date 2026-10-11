package db

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"rps/config"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QueryTracer struct {
	log *slog.Logger
}

var commentRegex = regexp.MustCompile(`(?m)^--.*$\n?`)

func (q *QueryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {

	// Format parameters cleanly
	formattedArgs := formatArgs(data.Args)
	cleanSQL := commentRegex.ReplaceAllString(data.SQL, "")
	q.log.Debug("Executing SQL Query",
		"sql", cleanSQL,
		"bindings", formattedArgs,
	)

	return ctx
}

func (q *QueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil {
		q.log.Error("SQL Query Failed", slog.Any("error", data.Err))
	}
}

// Helper to make pgtype and struct parameters readable in logs
func formatArgs(args []any) []string {
	result := make([]string, len(args))
	for i, arg := range args {
		result[i] = fmt.Sprintf("$%d = %v", i+1, arg)
	}
	return result
}

func NewPostgres(cfg config.DatabaseConfig, log *slog.Logger) (*pgxpool.Pool, error) {

	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = cfg.MaxOpenConns
	poolConfig.MinConns = cfg.MinIdleConns
	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(
		context.Background(),
		poolConfig,
	)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
