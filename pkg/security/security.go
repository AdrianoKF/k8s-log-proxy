package security

import (
	"fmt"
	"log"
	"strings"

	"github.com/adrianokf/k8s-log-proxy/pkg/config"
	"github.com/gobwas/glob"
)

func checkNamespace(namespace string, allowedNamespaces []string) bool {
	for _, ns := range allowedNamespaces {
		if strings.Contains(ns, "*") {
			// Allowed namespace contains a wildcard
			p := glob.MustCompile(ns)
			log.Printf("Glob %v, pattern %v\n", ns, p)
			if p.Match(namespace) {
				return true
			}
		} else if ns == namespace {
			return true
		}
	}
	return false
}

func CheckPermissions(namespace, podId string) (bool, error) {
	appConfig, err := config.ReadAppConfigFromFile()
	if err != nil {
		return false, err
	}

	// Only allow pods in allowlisted namespaces
	isAllowed := checkNamespace(namespace, appConfig.AllowedNamespaces)
	if isAllowed {
		return true, nil
	}
	return false, fmt.Errorf("forbidden: namespace '%v' not in %v", namespace, appConfig.AllowedNamespaces)
}
