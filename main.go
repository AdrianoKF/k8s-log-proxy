package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/adrianokf/k8s-log-proxy/pkg/config"
	"github.com/adrianokf/k8s-log-proxy/pkg/k8s"
	"github.com/adrianokf/k8s-log-proxy/pkg/logs"
	"github.com/adrianokf/k8s-log-proxy/pkg/security"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var lineReadLimit int64 = 8192

func LogError(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, msg)
	log.Println(msg)
}

func main() {
	kubeconfig, err := k8s.ReadKubeConfig()
	if err != nil {
		panic(err.Error())
	}
	client, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	config.Init()
	log.Printf("%+v\n", config.Config)

	// Main app
	http.HandleFunc("/", makeHandler(client))

	// Favicon
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "resources/favicon.ico")
	})

	// Kubernetes health check endpoints
	http.HandleFunc("/.healthz/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "OK")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func makeHandler(client kubernetes.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("URL: %s\n", r.URL.Path)

		target, err := logs.ParseUrl(r.URL.Path)
		if err != nil {
			return
		}

		if _, err := security.CheckPermissions(target.Namespace, target.Pod); err != nil {
			LogError(w, http.StatusForbidden, fmt.Sprintln("forbidden:", err.Error()))
			return
		}

		pod, err := client.CoreV1().Pods(target.Namespace).Get(context.TODO(), target.Pod, metav1.GetOptions{})
		if err != nil {
			LogError(w, http.StatusInternalServerError, fmt.Sprintln("error getting pod info:", err.Error()))
			return
		}

		var container string
		if target.Container != "" {
			container = target.Container
		} else {
			container = pod.Spec.Containers[0].Name
		}
		logOptions := &v1.PodLogOptions{
			Container:  container,
			Follow:     false,
			Timestamps: true,
			TailLines:  &lineReadLimit,
		}

		logs, err := logs.ReadLogs(client, pod.Namespace, pod.Name, logOptions)
		if err != nil {
			LogError(w, http.StatusInternalServerError, fmt.Sprintln("error reading logs:", err.Error()))
			return
		}
		fmt.Fprintln(w, logs)
	}
}
