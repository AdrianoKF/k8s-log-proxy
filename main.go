package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/adrianokf/k8s-log-proxy/pkg/k8s"
	"github.com/adrianokf/k8s-log-proxy/pkg/logs"
	"github.com/adrianokf/k8s-log-proxy/pkg/security"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var lineReadLimit int64 = 8192

func main() {
	config, err := k8s.ReadKubeConfig()
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

		_, err := security.CheckPermissions(namespace, podName)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			msg := fmt.Sprintf("forbidden: %s\n", err.Error())
			fmt.Fprintln(w, msg)
			log.Println(msg)
			return
		}

		pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			msg := fmt.Sprintf("Can't get pod info: %s\n", err.Error())
			fmt.Fprintln(w, msg)
			log.Println(msg)
			return
		}

		container := pod.Spec.Containers[0].Name
		logOptions := &v1.PodLogOptions{
			Container:  container,
			Follow:     false,
			Timestamps: true,
			TailLines:  &lineReadLimit,
		}

		logs, err := logs.ReadLogs(client, pod.Namespace, pod.Name, logOptions)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err.Error())
			log.Println(err.Error())
		}
		fmt.Fprintln(w, logs)
	}
}
