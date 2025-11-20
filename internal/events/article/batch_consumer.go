package article

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/jym0818/mywe/internal/repository"
	"github.com/jym0818/mywe/pkg/logger"
	"github.com/jym0818/mywe/pkg/saramax"
)

type BatchConsumerReadEvent struct {
	l      logger.Logger
	client sarama.Client
	repo   repository.InteractiveRepository
}

func NewBatchConsumerReadEvent(l logger.Logger, client sarama.Client) *BatchConsumerReadEvent {
	return &BatchConsumerReadEvent{
		l:      l,
		client: client,
	}
}

func (b *BatchConsumerReadEvent) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", b.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(), []string{"read_article"}, saramax.NewBatchHandler[ReadEvent](b.l, b.Consume, saramax.WithBatchDuration[ReadEvent](time.Second), saramax.WithBatchSize[ReadEvent](10)))
		if er != nil {
			//记录日志
		}
	}()
	return nil
}

func (b *BatchConsumerReadEvent) Consume(msgs []*sarama.ConsumerMessage, ts []ReadEvent) error {
	ids := make([]int64, 0, len(ts))
	bizs := make([]string, 0, len(ts))
	for _, evt := range ts {
		ids = append(ids, evt.Aid)
		bizs = append(bizs, "article")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := b.repo.BatchIncrReadCnt(ctx, bizs, ids)
	if err != nil {
		r.l.Error("批量增加阅读计数失败",
			logger.Field{Key: "ids", Value: ids},
			logger.Error(err))
	}
	return nil

}
