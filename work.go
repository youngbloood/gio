package gio

import (
	"context"

	"github.com/youngbloood/gio/log"
)

type Jober interface {
	Run(Result) error
}

type work struct {
	ctx      context.Context
	id       int             // work id
	workPool chan chan Jober // work pool
	jobChan  chan Jober      // worker 从pool中取出jober进行处理
	result   Result          // 结果集
}

func newWorker(workerPool chan chan Jober, result Result, id int) *work {
	return &work{
		id:       id,
		workPool: workerPool,
		jobChan:  make(chan Jober),
		result:   result,
		ctx:      context.Background(),
	}
}

func (w *work) start() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Handle(err)
			}
		}()

		for {
			//将worker的JobChannel放入master的workerPool中
			w.workPool <- w.jobChan
			select {
			//从JobChannel中获取Job进行处理，JobChannel是同步通道，会阻塞于此
			case job := <-w.jobChan:
				//处理这个job
				//并将处理得到的结果存入master中的结果集
				if err := job.Run(w.result); err != nil {
					log.Handle(err)
				}
			//停止信号
			case <-w.ctx.Done():
				return
			}
		}
	}()
}

func (w *work) stop() {
	go func() {
		var cancel context.CancelFunc
		w.ctx, cancel = context.WithCancel(w.ctx)
		cancel()
	}()
}
