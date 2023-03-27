package logs

import (
	"context"
	"io"
	"log"

	"github.com/adrianokf/k8s-log-proxy/pkg/security"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

func ReadLogs(client kubernetes.Interface, namespace, podId string, logOptions *v1.PodLogOptions) (string, error) {
	_, err := security.CheckPermissions(namespace, podId)
	if err != nil {
		panic(err.Error())
	}

	readCloser, err := client.CoreV1().Pods(namespace).GetLogs(podId, logOptions).Stream(context.TODO())

	if err != nil {
		log.Printf("Could not query pod logs: %s\n", err.Error())
		return "", err
	}
	defer readCloser.Close()

	result, err := io.ReadAll(readCloser)
	if err != nil {
		log.Printf("Could not read pod logs: %s\n", err.Error())
		return "ERROR", err
	}
	return string(result), nil
}
