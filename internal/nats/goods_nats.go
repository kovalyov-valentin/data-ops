package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kovalyov-valentin/data-ops/internal/domain/models"
	"github.com/kovalyov-valentin/data-ops/internal/storage"
	"github.com/nats-io/nats.go"
	"time"
)

type ItemRepo struct {
	stream    nats.JetStreamContext
	goodsRepo storage.Goods
}

func NewGoodsRepo(goodsRepo storage.Goods, js nats.JetStreamContext) (*ItemRepo, error) {
	cfg := &nats.StreamConfig{
		Name:      "EVENTS",
		Subjects:  []string{"events.>"},
		Retention: nats.WorkQueuePolicy,
	}

	_, err := js.AddStream(cfg)
	if err != nil {
		panic(err)
	}
	return &ItemRepo{
		goodsRepo: goodsRepo,
		stream:    js,
	}, nil
}

func (i *ItemRepo) UpdateItem(ctx context.Context, name, description string, id, projectsId int) (models.Goods, error) {
	var ce *models.ClickhouseEvent
	it, err := i.goodsRepo.UpdateGoods(ctx, name, description, id, projectsId)

	if err != nil {
		return it, err
	}

	ce.Id = it.ID
	ce.ProjectId = it.ProjectsID
	ce.Name = it.Name
	ce.Description = it.Description
	ce.Priority = it.Priority
	ce.Removed = it.Removed
	ce.EventTime = time.Now()
	ce.Description = it.Description
	i.sendEvent(ce)
	return it, err
}
func (i *ItemRepo) DeleteItem(ctx context.Context, id, projectsId int) (models.Goods, error) {
	it, err := i.goodsRepo.DeleteGoods(ctx, id, projectsId)

	ce := &models.ClickhouseEvent{
		Id:          id,
		ProjectId:   projectsId,
		EventTime:   time.Now(),
		Name:        it.Name,
		Description: it.Description,
		Priority:    it.Priority,
		Removed:     it.Removed,
	}
	i.sendEvent(ce)

	return it, err
}

func (i *ItemRepo) sendEvent(event *models.ClickhouseEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}
	s, err := i.stream.Publish("events.item", data)
	fmt.Println(s.Domain, s.Stream, s.Sequence, err, event)

}
