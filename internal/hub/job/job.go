package job

import (
	"time"

	"github.com/naturalselectionlabs/rss3-global-indexer/common/shedlock"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type Job interface {
	Name() string
	Spec() string
	Timeout() time.Duration
	Run() error
}

// Need to make Job compatible with cron.Job
var _ cron.Job = (*cronJob)(nil)

type cronJob struct {
	employer *shedlock.Employer
	job      Job
}

func (c *cronJob) Run() {
	zap.L().Info("worker job start", zap.String("name", c.job.Name()), zap.String("spec", c.job.Spec()), zap.Duration("timeout", c.job.Timeout()))

	if err := c.job.Run(); err != nil {
		logrus.Errorf("job %s throws an error: %s", c.job.Name(), err)
	}

	zap.L().Info("worker job end", zap.String("name", c.job.Name()), zap.String("spec", c.job.Spec()), zap.Duration("timeout", c.job.Timeout()))
}

func NewCronJob(employer *shedlock.Employer, job Job) cron.Job {
	return &cronJob{
		employer: employer,
		job:      job,
	}
}
