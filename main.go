package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var lineReadLimit int64 = 8192

func main() {
	config, err := readKubeConfig()
	if err != nil {
		panic(err.Error())
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/", makeHandler(client))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func readKubeConfig() (*rest.Config, error) {
	// Attempt to read in-cluster config
	config, _ := rest.InClusterConfig()
	if config != nil {
		return config, nil
	}

	// Attempt to build config from kubeconfig file
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func readLogs(client kubernetes.Interface, namespace, podID string, logOptions *v1.PodLogOptions) (string, error) {
	readCloser, err := client.CoreV1().Pods(namespace).GetLogs(podID, logOptions).Stream(context.TODO())

	if err != nil {
		log.Printf("Could not query pod logs: %s\n", err.Error())
		return err.Error(), nil
	}
	defer readCloser.Close()

	result, err := io.ReadAll(readCloser)
	if err != nil {
		log.Printf("Could not read pod logs: %s\n", err.Error())
		return "ERROR", err
	}
	return string(result), nil
}

func makeHandler(client kubernetes.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("URL: %s\n", r.URL.Path)
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/logs/"), "/")
		if len(parts) != 2 {
			return
		}
		namespace, podName := parts[0], parts[1]
		if namespace == "" || podName == "" {
			return
		}

		pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Can't get pod info: %s\n", err.Error())
			panic(err.Error())
		}

		container := pod.Spec.Containers[0].Name
		logOptions := &v1.PodLogOptions{
			Container:  container,
			Follow:     false,
			Timestamps: true,
			TailLines:  &lineReadLimit,
		}

		logs, err := readLogs(client, pod.Namespace, pod.Name, logOptions)
		if err != nil {
			panic(err.Error())
		}
		fmt.Fprintln(w, logs)
	}
}
