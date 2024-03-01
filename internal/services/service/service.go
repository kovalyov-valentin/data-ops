package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kovalyov-valentin/data-ops/internal/domain/models"
	"github.com/kovalyov-valentin/data-ops/internal/storage"
	"github.com/nats-io/nats.go"
	"time"
)

type EventSaver struct {
	repo storage.Event
	njs  nats.JetStreamContext
}

func NewService(repo storage.Event, njs nats.JetStreamContext) *EventSaver {
	es := &EventSaver{
		njs:  njs,
		repo: repo,
	}
	return es
}

func (es *EventSaver) CheckBatch(ctx context.Context) {

	for {
		sub, _ := es.njs.PullSubscribe("events.item",
			"worker",
			nats.PullMaxWaiting(128),
			nats.BindStream("EVENTS"),
		)
		if _, ok := ctx.Deadline(); ok {
			break
		}

		time.Sleep(time.Second * 10)
		msgs, err := sub.FetchBatch(100, nats.Context(ctx))
		if err != nil {
			continue
		}

		batch := make([]models.ClickhouseEvent, 0, 100)
		for msg := range msgs.Messages() {

			var ev models.ClickhouseEvent
			_ = json.Unmarshal(msg.Data, &ev)
			_ = msg.Ack()
			batch = append(batch, ev)
		}
		if len(batch) != 0 {
			err = es.repo.CreateEvent(ctx, batch)
		}
		fmt.Println(batch, err)
	}
}

func (es *EventSaver) Start(ctx context.Context) {
	go es.CheckBatch(ctx)
}
