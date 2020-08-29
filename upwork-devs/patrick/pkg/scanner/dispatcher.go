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

// File processing settings
type ProcessSettings struct {
	// A pool of workers channels that are registered with the dispatcher
	SourceFolder        string // Folder where source files are places
	ProcessingFolder    string // Folder where files are copied for processing, one folder per file
	OutputFolder        string // Processed files are places here before being exported to an object store (minio, s3, ...)
	ProcessPodImage     string // Image used to start the process pod
	ProcessPodNamespace string // Namespace where those pods will be created
	StorageAccessKey    string
	StorageSecretKey    string
	StorageBucket       string
	StorageEndpoint     string
}

//----------------------------------------------
// Exports
//----------------------------------------------
func NewDispatcher(maxWorkers int, kubeClient client.Client, processSettings *ProcessSettings) *Dispatcher {
	KubeClient = kubeClient
	Ps = processSettings
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
