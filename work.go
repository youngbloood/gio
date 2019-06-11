package gio

import (
	"context"
)

// Jober interface
type Jober interface {
	Run(Result) error
}

type work struct {
	ctx     context.Context
	id      int        // work id
	jobChan chan Jober // unbuffered channel. read jober and handle it

	pool *pool
}

func newWorker(id int, pool *pool) *work {
	return &work{
		id:      id,
		jobChan: make(chan Jober),
		ctx:     context.Background(),
		pool:    pool,
	}
}

func (w *work) start() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				Handle(err)
			}
		}()

		for {
			// write the w.jobChan into the w.pool.jobPool.
			w.pool.jobPool <- w.jobChan
			select {
			// read a job from w.jobChan, if there nothing, it will block.
			case job := <-w.jobChan:
				// handle the job.
				if err := job.Run(w.pool.result); err != nil {
					// push err to pool.
					w.pool.addErr(err)
					// log handle err.
					Handle(err)
					// not return. Do next jober.
				}
			// stop signal.
			case <-w.ctx.Done():
				return
			}
		}
	}()
}

// stop the work.
func (w *work) stop() {
	go func() {
		var cancel context.CancelFunc
		w.ctx, cancel = context.WithCancel(w.ctx)
		cancel()
	}()
}
