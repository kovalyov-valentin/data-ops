package storage

import (
	"context"
	"database/sql"
	"github.com/kovalyov-valentin/data-ops/internal/domain/models"
	"github.com/kovalyov-valentin/data-ops/internal/storage/postgres"
	rdb "github.com/kovalyov-valentin/data-ops/internal/storage/redis"
	"github.com/redis/go-redis/v9"
)

type Goods interface {
	CreateGoods(ctx context.Context, name string, projectsId int) (models.Goods, error)
	UpdateGoods(ctx context.Context, name, description string, id, projectsId int) (models.Goods, error)
	DeleteGoods(ctx context.Context, id, projectsId int) (models.Goods, error)
	GetGood(ctx context.Context, id, projectsId int) (models.Goods, error)
	GetGoods(ctx context.Context) ([]models.Goods, error)
}

type Cache interface {
	GetGoods(ctx context.Context) ([]models.Goods, error)
	GetGood(ctx context.Context, id, projectsId int) (models.Goods, error)
	Update(ctx context.Context, name, description string, id, projectsId int) (models.Goods, error)
	Delete(ctx context.Context, id, projectsId int) (models.Goods, error)
}

type Event interface {
	CreateEvent(ctx context.Context, event []models.ClickhouseEvent) error
}

type Repository struct {
	Goods Goods
	Cache Cache
	Event Event
}

func NewRepository(db *sql.DB, client *redis.Client) *Repository {
	goodsRepo := postgres.NewGoodsPostgres(db)
	redisGoods := rdb.NewGoodsRedis(client, goodsRepo)
	//natsGoods, _ := natsService.NewGoodsRepo(redisGoods, njs)
	//eventRepo := clickhouse.NewRepositoryClickhouse(db)

	return &Repository{
		Goods: goodsRepo,
		Cache: redisGoods,
		//Event: eventRepo,
	}
}
