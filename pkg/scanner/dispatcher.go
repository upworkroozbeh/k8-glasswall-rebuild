package scanner

import "sigs.k8s.io/controller-runtime/pkg/client"

//----------------------------------------------
// Imports
//----------------------------------------------

//----------------------------------------------
// Types
//----------------------------------------------
type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	maxWorkers int
	WorkerPool chan chan Job
}

//----------------------------------------------
// Exports
//----------------------------------------------
func NewDispatcher(maxWorkers int, kubeClient client.Client) *Dispatcher {
	KubeClient = kubeClient
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkerPool: pool, maxWorkers: maxWorkers}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

//----------------------------------------------
// Local Funcs
//----------------------------------------------
func (d *Dispatcher) dispatch() {
	for {

		select {
		case job := <-JobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}
