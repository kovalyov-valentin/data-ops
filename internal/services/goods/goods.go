package goods

import (
	"context"
	"errors"
	"fmt"
	"github.com/kovalyov-valentin/data-ops/internal/domain/models"
	"github.com/kovalyov-valentin/data-ops/internal/storage"
)

type UseCase struct {
	repo storage.Goods
}

func NewGoodsUseCase(repo storage.Goods) *UseCase {
	return &UseCase{repo: repo}
}

func (i *UseCase) CreateGoods(ctx context.Context, name string, projectsId int) (models.Goods, error) {
	if name == "" {
		return models.Goods{}, errors.New("failed to create goods")
	}
	return i.repo.CreateGoods(ctx, name, projectsId)
}

func (i *UseCase) UpdateGoods(ctx context.Context, name, description string, id, projectsId int) (models.Goods, error) {
	if name == "" {
		return models.Goods{}, errors.New("failed to update goods")
	}

	return i.repo.UpdateGoods(ctx, name, description, id, projectsId)
}

func (i *UseCase) DeleteGoods(ctx context.Context, id, projectsId int) (models.Goods, error) {
	return i.repo.DeleteGoods(ctx, id, projectsId)
}

func (i *UseCase) GetGood(ctx context.Context, id, projectsId int) (models.Goods, error) {
	return i.repo.GetGood(ctx, id, projectsId)
}

func (i *UseCase) GetGoods(ctx context.Context, limit, offset int) (models.GoodsResponse, error) {
	const op = "services.goods.GetGoods"
	allGoods, err := i.repo.GetGoods(ctx, limit, offset)
	if err != nil {
		return models.GoodsResponse{}, fmt.Errorf("%s: %w", op, err)
	}
	return allGoods, nil
}
