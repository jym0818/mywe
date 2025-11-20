package article

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/jym0818/mywe/internal/repository"
	"github.com/jym0818/mywe/pkg/logger"
	"github.com/jym0818/mywe/pkg/saramax"
)

type InteractiveReadEventConsumer struct {
	l      logger.Logger
	repo   repository.InteractiveRepository
	client sarama.Client
}

func NewInteractiveReadEventConsumer(l logger.Logger, repo repository.InteractiveRepository, client sarama.Client) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{
		l:      l,
		repo:   repo,
		client: client,
	}
}
func (i *InteractiveReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(), []string{"read_article"}, saramax.NewHandler[ReadEvent](i.l, i.Consume))
		if er != nil {
			//记录日志
		}
	}()
	return nil

}
func (i *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return i.repo.IncrReadCnt(ctx, "article", t.Aid)
}
