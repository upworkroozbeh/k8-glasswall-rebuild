package scanner

import (
	"context"
	"log"
	"strconv"
	"strings"

	kcorev1 "k8s.io/api/core/v1"
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// This method will run the scan again'st a file encapsulated the job.
// It starts a pod, passing it the path to the folder where the file is
// TODO : Copy that file to a location where the pod can access it. probably a shared location like NFS mount, etc
func (job Job) processFile(kubeClient client.Client) {

	log.Println("Processing file : " + job.File)

	// Env variable
	base_path_in_pod := "/tmp/files/"

	task_id_str := strconv.Itoa(job.TaskID)
	input_location := base_path_in_pod + "process-file-" + task_id_str
	input_location = strings.Replace(input_location, "/", "\\/", -1)

	pod := &kcorev1.Pod{
		ObjectMeta: kmetav1.ObjectMeta{
			Name:      "rebuild-" + task_id_str + "-" + job.Batch,
			Namespace: "test",
		},
		Spec: kcorev1.PodSpec{
			Containers:    []kcorev1.Container{{Name: "rebuild", Image: "azopat/gw-rebuild", Command: []string{"/bin/sh", "-c"}, Args: []string{"sed 's/INPUT_LOCATION/" + input_location + "/g' /home/glasswall/Config.ini.tmpl > /home/glasswall/Config.ini; glasswallCLI -config=/home/glasswall/Config.ini -xmlconfig=/home/glasswall/Config.xml && sleep infinity"}}},
			RestartPolicy: kcorev1.RestartPolicyNever,
		},
	}

	err := kubeClient.Create(context.Background(), pod)
	if err != nil {
		log.Println(err.Error())
	}

}
