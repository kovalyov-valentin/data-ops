package postgres

import (
	"database/sql"
	"fmt"
	"github.com/kovalyov-valentin/data-ops/internal/config"
	_ "github.com/lib/pq"
)

//type Storage struct {
//	db *sql.DB
//}

const (
	tableGoods = "goods"
)

func NewPostgresDB(cfg config.PostgresDB) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
