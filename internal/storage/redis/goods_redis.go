package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kovalyov-valentin/data-ops/internal/domain/models"
	"github.com/kovalyov-valentin/data-ops/internal/storage/postgres"
	"github.com/redis/go-redis/v9"
	"time"
)

type GoodsRedis struct {
	cache   *redis.Client
	goodsPG *postgres.GoodsPostgres
}

func NewGoodsRedis(cache *redis.Client, goodsPG *postgres.GoodsPostgres) *GoodsRedis {
	return &GoodsRedis{
		cache:   cache,
		goodsPG: goodsPG,
	}
}

func (r *GoodsRedis) GetGoods(ctx context.Context) ([]models.Goods, error) {
	val, err := r.cache.Get(ctx, "GetGoods").Bytes()
	if err != nil {
		goodsList, err := r.goodsPG.GetGoods(ctx)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(goodsList)
		if err != nil {
			return goodsList, nil
		}
		r.cache.SetNX(ctx, "GetGoods", data, time.Minute)
		return goodsList, nil
	}

	list := make([]models.Goods, 0)
	err = json.Unmarshal(val, &list)
	if err != nil {
		return nil, err
	}
	return list, err
}
func (r *GoodsRedis) GetGood(ctx context.Context, id, projectsId int) (models.Goods, error) {
	goodsKey := fmt.Sprintf("GetGoods-%d-%d", id, projectsId)
	val, err := r.cache.Get(ctx, goodsKey).Bytes()

	if err != nil {
		it, err := r.goodsPG.GetGood(ctx, id, projectsId)
		if err != nil {
			return models.Goods{}, err
		}

		data, err := json.Marshal(it)
		if err != nil {
			return it, nil
		}
		r.cache.SetNX(ctx, goodsKey, data, time.Minute)
		return it, nil
	}

	var goods models.Goods

	err = json.Unmarshal(val, &goods)
	if err != nil {
		return goods, err
	}
	return goods, err
}

func (r *GoodsRedis) Update(ctx context.Context, name, description string, id, projectsId int) (models.Goods, error) {
	r.deleteKey(ctx, id, projectsId)
	return r.goodsPG.UpdateGoods(ctx, name, description, id, projectsId)
}

func (r *GoodsRedis) Delete(ctx context.Context, id, projectsId int) (models.Goods, error) {
	r.deleteKey(ctx, id, projectsId)
	return r.goodsPG.DeleteGoods(ctx, id, projectsId)
}

func (r *GoodsRedis) deleteKey(ctx context.Context, id, projectsId int) {
	goodsKey := fmt.Sprintf("GetGood-%d-%d", id, projectsId)
	r.cache.Del(ctx, goodsKey)
	r.cache.Del(ctx, "GetGoods")
}
