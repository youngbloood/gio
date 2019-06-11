package gio

type master struct {
	workerPool chan chan Jober //worker池

	Result Result //存放worker处理后的结果集

	jobQueue chan Jober //待处理的任务chan

	workerList []*work //存放worker列表，用于停止worker

}

func NewMaster(maxWork int, result Result) *master {

	return &master{
		workerPool: make(chan chan Jober, maxWork),
		jobQueue:   make(chan Jober, 1.5*maxWork),
		Result:     result,
	}
}

// Run start all work
func (m *master) Run() {
	for i := 0; i <= len(m.workerPool); i++ {
		work := newWorker(m.workerPool, m.Result, i)
		m.workerList = append(m.workerList, work)
		work.start()
	}
	go m.dispatch()
}

func (m *master) dispatch() {
	for {
		select {
		case job := <-m.jobQueue:
			go func(job Jober) {
				//从workerPool中取出一个worker的JobChannel
				jobChan := <-m.workerPool
				//向这个JobChannel中发送job，worker中的接收配对操作会被唤醒
				jobChan <- job
			}(job)
		}
	}
}

func (m *master) Push(job Jober) {
	m.jobQueue <- job
}

func (m *master) Stop() {
	for _, wrk := range m.workerList {
		wrk.stop()
	}
}
