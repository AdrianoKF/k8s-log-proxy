package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/adrianokf/k8s-log-proxy/internal/k8s"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

func readLogs(client kubernetes.Interface, namespace, podID string, logOptions *v1.PodLogOptions) error {
	readCloser, err := client.CoreV1().Pods(namespace).GetLogs(podID, logOptions).Stream(context.TODO())

	if err != nil {
		log.Printf("Could not query pod logs: %s\n", err.Error())
		return err
	}
	defer readCloser.Close()

	scanner := bufio.NewScanner(readCloser)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return nil
}

func main() {
	os.Setenv("KUBECONFIG", "/home/adriano/.kube/kubeconfig-aai")

	config, err := k8s.ReadKubeConfig()
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	logOptions := v1.PodLogOptions{
		Follow: true,
	}
	_ = readLogs(client, "adrian", "planning-poker-planning-poker-prepper-c995c8645-mdtd4", &logOptions)
}
