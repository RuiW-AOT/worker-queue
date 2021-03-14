package workerq

import "sync"

type Dispatcher struct {
	Workers  []*Worker
	JobQueue chan Job
	Wg       *sync.WaitGroup
}

func NewDispatcher(workerNum int, wg *sync.WaitGroup) *Dispatcher {
	return &Dispatcher{
		Workers:  make([]*Worker, workerNum),
		JobQueue: make(chan Job),
		Wg:       wg,
	}
}

func (d *Dispatcher) Start() {
	for i := range d.Workers {
		d.Workers[i] = NewWorker(i, d.JobQueue, make(chan int), d.Wg)
		d.Workers[i].Start()
	}
}

func (d *Dispatcher) Submit(job Job) {
	d.JobQueue <- job
}

func (d *Dispatcher) ShutDown() {
	for _, w := range d.Workers {
		w.Quit <- 1
	}
}
