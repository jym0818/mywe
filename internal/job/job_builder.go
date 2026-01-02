package job

import (
	"github.com/jym0818/mywe/pkg/logger"
	"github.com/robfig/cron/v3"
)

type cronJobFuncAdapter func() error

func (c cronJobFuncAdapter) Run() {
	_ = c()
}

type CronJobBuilder struct {
	l logger.Logger
}

func NewCronJobBuilder(l logger.Logger) *CronJobBuilder {
	return &CronJobBuilder{l: l}
}
func (b *CronJobBuilder) Build(job Job) cron.Job {
	//name := job.Name()
	return cronJobFuncAdapter(func() error {
		err := job.Run()
		if err != nil {
			//记录日志
		}
		return nil
	})
}
