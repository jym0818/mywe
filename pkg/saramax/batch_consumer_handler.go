package saramax

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/jym0818/mywe/pkg/logger"
)

type BatchHandler[T any] struct {
	l             logger.Logger
	fn            func(msgs []*sarama.ConsumerMessage, ts []T) error
	batchSize     int
	batchDuration time.Duration
}

func (b BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgsCh := claim.Messages()
	for {
		ctx, cancel := context.WithTimeout(context.Background(), b.batchDuration)
		msgs := make([]*sarama.ConsumerMessage, 0, b.batchSize)
		ts := make([]T, 0, b.batchSize)
		for i := 0; i < b.batchSize; i++ {
			done := false
			select {
			case <-ctx.Done():
				done = true
			case msg, ok := <-msgsCh:
				if !ok {
					cancel()
					return nil
				}
				var t T
				if err := json.Unmarshal(msg.Value, &t); err != nil {
					//记录日志
					continue
				}
				msgs = append(msgs, msg)
				ts = append(ts, t)

			}
			if done {
				break
			}
		}
		cancel()
		if len(msgs) == 0 {
			continue
		}
		err := b.fn(msgs, ts)
		if err != nil {
			b.l.Error("调用业务批量接口失败", logger.Error(err))

		}
		for _, msg := range msgs {
			session.MarkMessage(msg, "")
		}
	}

}

type BatchHandlerOption[T any] func(*BatchHandler[T])

func WithBatchSize[T any](size int) BatchHandlerOption[T] {
	return func(b *BatchHandler[T]) {
		if size > 0 {
			b.batchSize = size
		}
	}
}

func WithBatchDuration[T any](duration time.Duration) BatchHandlerOption[T] {
	return func(b *BatchHandler[T]) {
		if duration > 0 {
			b.batchDuration = duration
		}
	}
}

func NewBatchHandler[T any](l logger.Logger, fn func(msgs []*sarama.ConsumerMessage, ts []T) error, opts ...BatchHandlerOption[T]) *BatchHandler[T] {
	handler := &BatchHandler[T]{
		l:             l,
		fn:            fn,
		batchSize:     10,
		batchDuration: time.Second,
	}
	for _, opt := range opts {
		opt(handler)
	}
	return handler
}
