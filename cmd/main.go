package main

import (
	"fmt"

	"github.com/adrianokf/k8s-log-proxy/pkg/k8s"
	"github.com/adrianokf/k8s-log-proxy/pkg/logs"
	v1 "k8s.io/api/core/v1"
)

func main() {
	client, err := k8s.MakeK8sClient()
	if err != nil {
		panic(err.Error())
	}
	logOptions := v1.PodLogOptions{}
	logString, err := logs.ReadLogs(client, "istio-system", "istiod-958644599-lbqdx", &logOptions)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(logString)
}
