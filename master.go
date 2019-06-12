package gio

import (
	"context"
)

// Master .
type Master struct {
	pool     *pool
	jobQueue chan Jober // undo jober
}

// NewMaster return a new master.
func NewMaster(maxWork int, result Saver) *Master {
	return &Master{
		pool: &pool{
			ctx:     context.Background(),
			jobPool: make(chan chan Jober, maxWork),
			result:  result,
			errs:    make([]error, 0),
		},
		jobQueue: make(chan Jober, 2*maxWork),
	}
}

// Run start all work.
func (m *Master) Run() {
	for i := 0; i < cap(m.pool.jobPool); i++ {
		newWorker(i, m.pool).start()
	}
	go m.dispatch()
}

// dispath
func (m *Master) dispatch() {
	for {
		select {
		case job := <-m.jobQueue:
			go func(job Jober) {
				select {
				// read a jobchan from m.pool.jobPool.
				case jobChan := <-m.pool.jobPool:
					select {
					// write the job into jobchan.
					case jobChan <- job:
					case <-m.pool.ctx.Done():
						return
					}
				// if the master called Stop(), will do next and exit the goroutine.
				// else if not called, it will blocked until read a jobchan from m.pool.jobPool(behind the work had handle the job then push itself into m.pool.jobPool).
				case <-m.pool.ctx.Done():
					return
				}

			}(job)
		}
	}
}

// Push jobs into jobqueue.
func (m *Master) Push(jobs ...Jober) {
	for _, job := range jobs {
		m.jobQueue <- job
	}
}

// Stop . stop all the work.
func (m *Master) Stop() {
	var cancel context.CancelFunc
	m.pool.ctx, cancel = context.WithCancel(m.pool.ctx)
	cancel()
}

// GetErrs return all errs.
func (m *Master) GetErrs() []error {
	return m.pool.getErrs()
}

// GetResult return result
// func (m *Master) GetResult() Result {
// 	return m.pool.result
// }
