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
	k8sclient, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		log.Println("config not found, exiting")
		os.Exit(1)
	}

	// This should be injected as env variables, represent amount of simultaneaus job
	maxQueue := 5
	maxWorker := 5

	// Shared volume folder root is /data and in contains a folder src-files where the source files should be copied
	sourceFolder := "/data/src-files"
	processingFolder := "/data/processing-files"
	outputFolder := "/data/processed-files"
	image := "azopat/gw-rebuild"
	namespace := "test"
	storageAccessKey := "minio"
	storageSecretKey := "minio123"
	storageBucket := "glasswall"
	storageEndpoint := "http://minio.default.svc.cluster.local:9000"

	processSettings := &scanner.ProcessSettings{SourceFolder: sourceFolder, ProcessingFolder: processingFolder, ProcessPodImage: image, ProcessPodNamespace: namespace, OutputFolder: outputFolder, StorageAccessKey: storageAccessKey, StorageSecretKey: storageSecretKey, StorageBucket: storageBucket, StorageEndpoint: storageEndpoint}

	// Workers initialization
	scanner.JobQueue = make(chan scanner.Job, maxQueue)
	dispatcher := scanner.NewDispatcher(maxWorker, k8sclient, processSettings)
	dispatcher.Run()

	// Starting a scan
	scanProcessor := scanner.ScanProcessor{Folder: sourceFolder, Batch: uuid.New().String(), ContainerImage: image, Namespace: namespace}
	scanProcessor.ScanFiles()

	// This is just to keep the pod running for now.
	for {

	}

}
