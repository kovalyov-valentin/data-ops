package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kovalyov-valentin/data-ops/internal/config"
	_ "github.com/mailru/go-clickhouse/v2"
)

func NewClickhouseDB(ctx context.Context, cfg config.Clickhouse) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("http://%s:%s/%s", cfg.Host, cfg.HttpPort, cfg.DB)
	db, err := sql.Open("chhttp", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error to connection postgres: %v", err)
	}
	go func(ctx context.Context) {
		<-ctx.Done()
		db.Close()
	}(ctx)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("connection test error: %w", err)
	}

	return db, nil
}
