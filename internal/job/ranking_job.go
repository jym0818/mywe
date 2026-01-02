package job

import (
	"context"
	"time"

	"github.com/jym0818/mywe/internal/service"
)

type RankingJob struct {
	svc     service.RankingService
	timeout time.Duration
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	_, err := r.svc.TopN(ctx)
	return err
}

func NewRankingJob(svc service.RankingService, timeout time.Duration) *RankingJob {
	return &RankingJob{
		svc:     svc,
		timeout: timeout,
	}
}
