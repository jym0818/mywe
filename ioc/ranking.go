package ioc

import (
	"time"

	"github.com/jym0818/mywe/internal/job"
	"github.com/jym0818/mywe/internal/service"
	"github.com/jym0818/mywe/pkg/logger"
	"github.com/robfig/cron/v3"
)

func InitRankingJob(svc service.RankingService) *job.RankingJob {
	return job.NewRankingJob(svc, time.Second*30)
}

func InitJobs(l logger.Logger, rankingJob *job.RankingJob) *cron.Cron {
	res := cron.New(cron.WithSeconds())
	cbd := job.NewCronJobBuilder(l)
	_, err := res.AddJob("0 */3 * * * ?", cbd.Build(rankingJob))
	if err != nil {
		panic(err)
	}
	return res
}
