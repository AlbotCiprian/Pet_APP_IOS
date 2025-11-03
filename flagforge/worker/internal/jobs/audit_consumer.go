package jobs

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type AuditConsumer struct {
	rdb *redis.Client
	log zerolog.Logger
}

func NewAuditConsumer(rdb *redis.Client, log zerolog.Logger) *AuditConsumer {
	return &AuditConsumer{rdb: rdb, log: log}
}

func (c *AuditConsumer) Run(ctx context.Context) {
	sub := c.rdb.Subscribe(ctx, "flagforge.invalidate")
	defer sub.Close()

	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			c.log.Error().Err(err).Msg("subscription error")
			continue
		}

		c.log.Info().Str("project_id", msg.Payload).Msg("received invalidation")
		// Simulate cache rebuild latency
		time.Sleep(100 * time.Millisecond)
	}
}
