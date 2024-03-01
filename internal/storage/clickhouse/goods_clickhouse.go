package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kovalyov-valentin/data-ops/internal/domain/models"
	"strings"
)

type RepositoryClickhouse struct {
	db *sql.DB
}

func NewRepositoryClickhouse(db *sql.DB) *RepositoryClickhouse {
	return &RepositoryClickhouse{db: db}
}

func (r *RepositoryClickhouse) CreateEvent(ctx context.Context, clickhouseEvents []models.ClickhouseEvent) error {
	query := "INSERT INTO goods (Id, CampaignId, Name, Description, Priority, Removed, EventTime) VALUES "
	var valueStrings []string
	var valueArgs []interface{}

	for _, ce := range clickhouseEvents {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, ce.Id, ce.ProjectId, ce.Name, ce.Description, ce.Priority, ce.Removed, ce.EventTime)
	}

	query += fmt.Sprintf("%s", strings.Join(valueStrings, ","))
	_, err := r.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("clickhouse.CreateEvent Exec %w", err)
	}
	return nil
}
