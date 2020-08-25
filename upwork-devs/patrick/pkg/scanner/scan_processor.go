package scanner

import (
	"io/ioutil"
	"log"
)

type ScanProcessor struct {
	Batch          string
	Folder         string
	ContainerImage string // Image containing scan tool - https://github.com/filetrust/Glasswall-Rebuild-SDK-Evaluation
	Namespace      string // the namespace where the pods will be created
}

func (s *ScanProcessor) ScanFiles() {

	log.Println("Scan processor on folder " + s.Folder)
	files, err := ioutil.ReadDir(s.Folder)
	if err != nil {
		log.Println(err.Error())
	}
	i := 1
	for _, f := range files {
		if !f.IsDir() {
			log.Println("File found : " + f.Name())
			job := Job{File: f.Name(), TaskID: i, Batch: s.Batch, ContainerImage: s.ContainerImage, Namespace: s.Namespace}
			i = i + 1
			JobQueue <- job
		}
	}

}
