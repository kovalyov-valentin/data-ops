package nats

import (
	"context"
	"fmt"
	"github.com/kovalyov-valentin/data-ops/internal/config"
	"github.com/nats-io/nats.go"
)

func NewNats(ctx context.Context, cfg config.Nats) (*nats.Conn, error) {
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%s", cfg.Host, cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect nats: %w", err)
	}

	go func(ctx context.Context) {
		<-ctx.Done()
		nc.Close()
	}(ctx)

	return nc, err

}
