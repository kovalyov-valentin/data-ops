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

func (p *GoodsPostgres) CreateGoods(ctx context.Context, name string, projectsId int) (models.Goods, error) {
	const op = "storage.postgres.CreateItem"

	query := `
		INSERT INTO 
		    goods (name, projects_id)
		VALUES 
		    ($1, $2)
		RETURNING id, projects_id, name, description, priority, removed, created_at
	`

	var goods models.Goods
	err := p.db.QueryRowContext(ctx, query, name, projectsId).
		Scan(&goods.ID, &goods.ProjectsID, &goods.Name, &goods.Description, &goods.Priority, &goods.Removed, &goods.CreatedAt)
	if err != nil {
		return models.Goods{}, fmt.Errorf("%s: %w", op, err)
	}

	return goods, nil
}

func (p *GoodsPostgres) UpdateGoods(ctx context.Context, name, description string, id, projectsId int) (models.Goods, error) {
	const op = "storage.postgres.UpdateGoods"

	tx, err := p.db.Begin()
	if err != nil {
		return models.Goods{}, fmt.Errorf("%s: %w", op, err)
	}

	queryUpdate := `
		UPDATE goods
		SET name = $1, description = $2
		WHERE id = $3 AND projects_id = $4
		RETURNING id
	`

	row := tx.QueryRowContext(ctx, queryUpdate, name, description, id, projectsId)
	if err := row.Err(); err != nil {
		return models.Goods{}, fmt.Errorf("%s: %w", op, err)
	}

	var idGoods int
	if err := row.Scan(&idGoods); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Goods{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	querySelect := `
		SELECT id, projects_id, name, description, priority, removed, created_at 
		FROM goods 
		WHERE id = $1`

	row = tx.QueryRowContext(ctx, querySelect, idGoods)
	if err := row.Err(); err != nil {
		return models.Goods{}, fmt.Errorf("%s: %w", op, err)
	}

	var goods models.Goods
	if err := row.Scan(
		&goods.ID,
		&goods.ProjectsID,
		&goods.Name,
		&goods.Description,
		&goods.Priority,
		&goods.Removed,
		&goods.CreatedAt,
	); err != nil {
		return models.Goods{}, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return models.Goods{}, fmt.Errorf("%s: %w", op, err)
	}

	return goods, nil

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
	const op = "storage.postgres.GetGood"

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Goods{}, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	query := `
		SELECT id, projects_id, name, description, priority, removed, created_at 
		FROM goods 
		WHERE id = $1 AND projects_id = $2
		`

	row := tx.QueryRowContext(ctx, query, id, projectsId)
	if err != nil {
		return models.Goods{}, fmt.Errorf("%s: %w", op, err)
	}

	var goods models.Goods

	err = row.Scan(
		&goods.ID,
		&goods.ProjectsID,
		&goods.Name,
		&goods.Description,
		&goods.Priority,
		&goods.Removed,
		&goods.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Goods{}, fmt.Errorf("%s: %w", op, err)
		}
		return models.Goods{}, fmt.Errorf("%s: %w", op, err)
	}

	if err = row.Err(); err != nil {
		return models.Goods{}, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return models.Goods{}, fmt.Errorf("%s: %w", op, err)
	}

	return goods, nil
}

func (p *GoodsPostgres) GetGoods(ctx context.Context, limit, offset int) (models.GoodsResponse, error) {
	// TODO: refactor
	const op = "storage.postgres.GetGoods"

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return models.GoodsResponse{}, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	query := `
		SELECT * 
		FROM goods 
		LIMIT $1 OFFSET $2
         `

	rows, err := tx.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return models.GoodsResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	var goods []models.Goods
	for rows.Next() {
		var good models.Goods
		err = rows.Scan(
			&good.ID,
			&good.ProjectsID,
			&good.Name,
			&good.Description,
			&good.Priority,
			&good.Removed,
			&good.CreatedAt,
		)
		if err != nil {
			return models.GoodsResponse{}, fmt.Errorf("%s: %w", op, err)
		}
		goods = append(goods, good)
	}

	if err = rows.Err(); err != nil {
		return models.GoodsResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	queryTotal := "SELECT COUNT(id) FROM goods"
	rowTotal := p.db.QueryRowContext(ctx, queryTotal)
	if err = rowTotal.Err(); err != nil {
		return models.GoodsResponse{}, fmt.Errorf("%s: %w", op, err)
	}
	var total int
	if err = rowTotal.Scan(&total); err != nil {
		return models.GoodsResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	queryRemoved := "SELECT COUNT(id) FROM goods WHERE removed = true"
	rowRemoved := p.db.QueryRowContext(ctx, queryRemoved)
	if err = rowRemoved.Err(); err != nil {
		return models.GoodsResponse{}, fmt.Errorf("%s: %w", op, err)
	}
	var removed int
	if err = rowRemoved.Scan(&removed); err != nil {
		return models.GoodsResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return models.GoodsResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	resp := models.GoodsResponse{
		Meta: models.Meta{
			Total:   total,
			Removed: removed,
			Limit:   limit,
			Offset:  offset,
		},
		Goods: goods,
	}

	return resp, nil
}
