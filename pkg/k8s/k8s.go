package k8s

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// ReadKubeConfig reads the kubeconfig file from either the in-cluster config or the local config
func ReadKubeConfig() (*rest.Config, error) {
	// Attempt to read in-cluster config
	config, _ := rest.InClusterConfig()
	if config != nil {
		return config, nil
	}

	// Attempt to build config from kubeconfig file
	var kubeconfig *string

	if val, found := os.LookupEnv("KUBECONFIG"); found {
		kubeconfig = &val
	} else {
		if home := homedir.HomeDir(); home != "" {
			val := filepath.Join(home, ".kube", "config")
			kubeconfig = &val
		}
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// MakeK8sClient creates a Kubernetes client for the current environment
func MakeK8sClient() (*kubernetes.Clientset, error) {
	config, err := ReadKubeConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
