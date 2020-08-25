package main

import (
	"log"
	"os"

	"github.com/google/uuid"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/filetrust/Open-Source/upwork/project-k8-glasswall-rebuild/pkg/scanner"
)

const ()

func main() {

	// Get a kubernetes client
	cl, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		log.Println("config not found, exiting")
		os.Exit(1)
	}

	// This should be injected as env variables, represent amount of simultaneaus job
	maxQueue := 5
	maxWorker := 5
	sourceFolder := "/tmp/files"
	image := "azopat/gw-rebuild"
	namespace := "test"

	scanner.JobQueue = make(chan scanner.Job, maxQueue)
	dispatcher := scanner.NewDispatcher(maxWorker, cl)
	dispatcher.Run()

	scanProcessor := scanner.ScanProcessor{Folder: sourceFolder, Batch: uuid.New().String(), ContainerImage: image, Namespace: namespace}
	scanProcessor.ScanFiles()

	// This is just to keep the pod running for now.
	for {

	}

}
