package config

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/adrianokf/k8s-log-proxy/pkg/k8s"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AppConfig struct {
	AllowedNamespaces []string `json:"allowedNamespaces"`
}

// ReadAppConfigFromK8s reads and parses the application config from a Kubernetes configmap
func ReadAppConfigFromK8s(namespace, configMap string) (*AppConfig, error) {
	client, err := k8s.MakeK8sClient()
	if err != nil {
		return nil, err
	}

	appConfig := AppConfig{
		AllowedNamespaces: make([]string, 0),
	}

	cm, err := client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMap, v1.GetOptions{})
	if err != nil {
		log.Printf("Application ConfigMap '%v' not accessible in namespace '%v', using defaults\n", configMap, namespace)
		return &appConfig, nil
	}

	err = json.Unmarshal([]byte(cm.Data["config.json"]), &appConfig)
	if err != nil {
		return nil, err
	}
	return &appConfig, nil
}

const DEFAULT_CONFIG_FILE = "/config.json"

// ReadAppConfigFromFile reads and parses the application config from a config file
func ReadAppConfigFromFile() (*AppConfig, error) {
	appConfig := AppConfig{
		AllowedNamespaces: make([]string, 0),
	}

	data, err := ioutil.ReadFile(DEFAULT_CONFIG_FILE)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &appConfig)
	if err != nil {
		return nil, err
	}
	return &appConfig, nil
}
