package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/kovalyov-valentin/data-ops/internal/domain/models"
)

type GoodsPostgres struct {
	db *sql.DB
}

func NewGoodsPostgres(db *sql.DB) *GoodsPostgres {
	return &GoodsPostgres{
		db: db,
	}
}

// TODO: transaction
// TODO: implement pagination
func (p *GoodsPostgres) CreateGoods(ctx context.Context, name string, projectsId int) (models.Goods, error) {
	const op = "postgres.CreateItem"

	query := `
		INSERT INTO goods (name, projects_id)
		VALUES ($1, $2)
		RETURNING id, created_at, priority
	`

	var item models.Goods
	err := p.db.QueryRowContext(ctx, query, name, projectsId).
		Scan(&item.ID, &item.CreatedAt, &item.Priority)
	if err != nil {
		// Обработка ошибки.
		return models.Goods{}, fmt.Errorf("%s: %w", op, err)
	}

	return item, nil
}

func (p *GoodsPostgres) UpdateGoods(ctx context.Context, name, description string, id, projectsId int) (models.Goods, error) {
	const op = "postgres.UpdateGoods"

	sqlQuery := `
		UPDATE goods
		SET name = $1, description = $2
		WHERE id = $3 AND projects_id = $4
		RETURNING id, created_at, priority
	`

	var item models.Goods
	err := p.db.QueryRowContext(ctx, sqlQuery, name, description, id, projectsId).Scan(&item.ID, &item.CreatedAt, &item.Priority)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return models.Goods{}, errors.New("not found goods")
		}
		return models.Goods{}, fmt.Errorf("%s: %w", op, err)
	}

	return item, nil

}

func (p *GoodsPostgres) DeleteGoods(ctx context.Context, id, projectsId int) (models.Goods, error) {
	const op = "postgres.DeleteGoods"

	query := `
		UPDATE goods
		SET removed = true
		WHERE id = $1 AND projects_id = $2
		RETURNING id, projects_id, name, description, priority, removed, created_at
	`

	var goods models.Goods
	err := p.db.QueryRowContext(ctx, query, id, projectsId).
		Scan(&goods.ID, &goods.ProjectsID, &goods.Name, &goods.Description, &goods.Priority, &goods.Removed, &goods.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Goods{}, fmt.Errorf("%s %w", op, errors.New("not found item"))
		}
		return models.Goods{}, fmt.Errorf("%s: %w", op, err)
	}

	return goods, nil

}

func (p *GoodsPostgres) GetGood(ctx context.Context, id, projectsId int) (models.Goods, error) {
	const op = "postgres.GetGood"

	var goods models.Goods

	query := `
		SELECT id, projects_id, name, description, priority, removed, created_at 
		FROM goods 
		WHERE id = $1 AND projects_id = $2`
	row := p.db.QueryRowContext(ctx, query, id, projectsId)

	err := row.Scan(
		&goods.ID,
		&goods.ProjectsID,
		&goods.Name,
		&goods.Description,
		&goods.Priority,
		&goods.Removed,
		&goods.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Goods{}, fmt.Errorf("%s: no rows found", op)
		}
		return models.Goods{}, fmt.Errorf("%s: %w", op, err)
	}

	return goods, nil
}

func (p *GoodsPostgres) GetGoods(ctx context.Context) ([]models.Goods, error) {
	const op = "postgres.GetCampaigns"

	itemList := make([]models.Goods, 0)

	sqlReq := "SELECT * FROM " + tableGoods

	rows, err := p.db.QueryContext(ctx, sqlReq)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var goods models.Goods
		err = rows.Scan(&goods.ID,
			&goods.ProjectsID,
			&goods.Name,
			&goods.Description,
			&goods.Priority,
			&goods.Removed,
			&goods.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		itemList = append(itemList, goods)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s rows.Err() %w", op, err)
	}

	return itemList, nil
}
