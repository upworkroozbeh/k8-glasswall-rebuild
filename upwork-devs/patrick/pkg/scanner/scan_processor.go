package scanner

import (
	"io/ioutil"
	"log"
)

type ScanProcessor struct {
	Batch  string
	Folder string
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
			job := Job{File: f.Name(), TaskID: i, Batch: s.Batch}
			i = i + 1
			JobQueue <- job
		}
	}

}
