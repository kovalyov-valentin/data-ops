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
	GetGoods(ctx context.Context, limit, offset int) (models.GoodsResponse, error)
}

type Cache interface {
	GetGoods(ctx context.Context) ([]models.GoodsResponse, error)
	GetGood(ctx context.Context, id, projectsId int) (models.Goods, error)
	Update(ctx context.Context, name, description string, id, projectsId int) (models.Goods, error)
	Delete(ctx context.Context, id, projectsId int) (models.Goods, error)
}

type repositoryClickhouse interface{}

type Queue interface {
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
	Cache      Cache
	Queue      Queue
}

func NewService(repos *storage.Repository) *Service {
	return &Service{
		Goods: goods.NewGoodsUseCase(repos.Goods),
	}
}
