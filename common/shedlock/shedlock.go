package shedlock

import (
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/robfig/cron/v3"
)

var (
	globalLocker sync.RWMutex
	globalRs     *redsync.Redsync
)

func GlobalRs() *redsync.Redsync {
	globalLocker.RLock()

	defer globalLocker.RUnlock()

	return globalRs
}

func ReplaceGlobalRs(rs *redsync.Redsync) {
	globalLocker.Lock()

	defer globalLocker.Unlock()

	globalRs = rs
}

type Employer struct {
	crontab *cron.Cron
}

func (e *Employer) AddJob(name, spec string, timeout time.Duration, cmd cron.Job) error {
	_, err := e.crontab.AddFunc(spec, func() {
		mutex := GlobalRs().NewMutex(name, redsync.WithExpiry(timeout))

		if err := mutex.Lock(); err != nil {
			return
		}

		defer func() {
			cmd.Run()

			_, _ = mutex.Unlock()
		}()
	})

	return err
}

func (e *Employer) Start() {
	e.crontab.Start()
}

func (e *Employer) Stop() {
	e.crontab.Stop()
}

func New() *Employer {
	return &Employer{
		crontab: cron.New(cron.WithSeconds()),
	}
}
