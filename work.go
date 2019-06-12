package gio

import (
	"github.com/youngbloood/gio/log"
)

// Jober interface
type Jober interface {
	Run(Saver) error
}

type work struct {
	id      int        // work id
	jobChan chan Jober // unbuffered channel. read jober and handle it

	pool *pool
}

func newWorker(id int, pool *pool) *work {
	return &work{
		id:      id,
		jobChan: make(chan Jober),
		pool:    pool,
	}
}

func (w *work) start() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if errv, ok := err.(error); ok {
					// push err to pool.
					w.pool.addErr(errv)
				}
				// log handle err.
				log.Handle(err)
			}
		}()
		for {
			select {
			// write the w.jobChan into the w.pool.jobPool.
			case w.pool.jobPool <- w.jobChan:
				select {
				// read a job from w.jobChan, if there nothing, it will block.
				case job := <-w.jobChan:
					// handle the job.
					if err := job.Run(w.pool.result); err != nil {
						// push err to pool.
						w.pool.addErr(err)
						// log handle err.
						log.Handle(err)
						// not return. Do next jober.
					}
				// stop signal.
				case <-w.pool.ctx.Done():
					return
				}
			// stop signal.
			case <-w.pool.ctx.Done():
				return
			}
		}
	}()
}
