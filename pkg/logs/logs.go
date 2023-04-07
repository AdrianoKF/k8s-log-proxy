package logs

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/adrianokf/k8s-log-proxy/pkg/security"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type LogTarget struct {
	Namespace string
	Pod       string
	Container string
}

// / ParseUrl parses a request URL path and returns its target (identified by namespace, pod, and optional container name)
func ParseUrl(urlPath string) (target LogTarget, err error) {
	parts := strings.Split(strings.TrimPrefix(urlPath, "/logs/"), "/")
	if len(parts) == 2 {
		target.Namespace, target.Pod = parts[0], parts[1]
	} else if len(parts) == 3 {
		target.Namespace, target.Pod, target.Container = parts[0], parts[1], parts[2]
	}
	if target.Namespace == "" || target.Pod == "" {
		err = fmt.Errorf("namespace or pod name empty: %+v", target)
	}
	return
}

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
