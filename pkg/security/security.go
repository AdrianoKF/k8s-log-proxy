package security

import (
	"fmt"

	"github.com/adrianokf/k8s-log-proxy/pkg/config"
)

func CheckPermissions(namespace, podId string) (bool, error) {
	appConfig, err := config.ReadAppConfigFromFile()
	if err != nil {
		return false, err
	}

	// Only allow pods in allowlisted namespaces
	for _, ns := range appConfig.AllowedNamespaces {
		if ns == namespace {
			return true, nil
		}
	}
	return false, fmt.Errorf("forbidden: namespace '%v' not in %v", namespace, appConfig.AllowedNamespaces)
}
