package sqlite

import (
	"app-dev/internal/domain"
	"app-dev/internal/driver/sqlcommon"
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func NewClient(ctx context.Context, path string) (domain.Client, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping sqlite database: %w", err)
	}

	return &sqlcommon.SQLClient{DB: db}, nil
}
