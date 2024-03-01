package services

import (
	"context"
	"github.com/kovalyov-valentin/data-ops/internal/domain/models"
	"github.com/kovalyov-valentin/data-ops/internal/services/goods"
	"github.com/kovalyov-valentin/data-ops/internal/storage"
)

type Goods interface {
	CreateGoods(ctx context.Context, name string, projectsId int) (models.Goods, error)
	UpdateGoods(ctx context.Context, name, description string, id, projectsId int) (models.Goods, error)
	DeleteGoods(ctx context.Context, id, projectsId int) (models.Goods, error)
	GetGood(ctx context.Context, id, projectsId int) (models.Goods, error)
	GetGoods(ctx context.Context) ([]models.Goods, error)
}

type EventSaver interface {
	Start(ctx context.Context)
}

type Event interface {
	CreateEvent(ctx context.Context, event []models.ClickhouseEvent) error
}

type Service struct {
	Goods      Goods
	EventSaver EventSaver
}

func NewService(repos *storage.Repository) *Service {
	return &Service{
		Goods: goods.NewGoodsUseCase(repos.Goods),
	}
}
