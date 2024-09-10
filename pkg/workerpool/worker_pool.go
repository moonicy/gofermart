package workerpool

type WorkerPool struct {
	workerCount int
	chJob       chan Job
	worker Worker
}

type Worker func(<-chan Job)

type Job func() error

func NewWorkerPool(w Worker, workerCount int) *WorkerPool {
	chJobs := make(chan Job)
	return &WorkerPool{workerCount: workerCount, worker: w, chJob: chJobs}
}

func (wp *WorkerPool) AddJob(job Job) {
	wp.chJob <- job
}

func (wp *WorkerPool) Run() {
	for i := 0; i < wp.workerCount; i++ {
		go wp.worker(wp.chJob)

	}
}

func (wp *WorkerPool) Close() {
	close(wp.chJob)
}
