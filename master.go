package gio

// Master .
type Master struct {
	pool *pool

	jobQueue   chan Jober // undo jober
	workerList []*work    // work list
}

// NewMaster return a new master.
func NewMaster(maxWork int, result Result) *Master {

	return &Master{
		pool: &pool{
			jobPool: make(chan chan Jober, maxWork),
			result:  result,
			errs:    make([]error, 0),
		},
		jobQueue:   make(chan Jober, 2*maxWork),
		workerList: make([]*work, maxWork),
	}
}

// Run start all work.
func (m *Master) Run() {
	for i := 0; i < len(m.pool.jobPool); i++ {
		work := newWorker(i, m.pool)
		m.workerList[i] = work
		work.start()
	}
	go m.dispatch()
}

func (m *Master) dispatch() {
	for {
		select {
		case job := <-m.jobQueue:
			go func(job Jober) {
				// read a jobchan from m.pool.jobPoll.
				jobChan := <-m.pool.jobPool
				// write the job into jobchan.
				jobChan <- job
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
	for _, wrk := range m.workerList {
		wrk.stop()
	}
}

// GetErrs return all errs.
func (m *Master) GetErrs() []error {
	return m.pool.getErrs()
}
