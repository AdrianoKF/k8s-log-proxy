package security

import (
	"fmt"
	"testing"
)

func TestCheckNamespace(t *testing.T) {
	tests := []struct {
		namespace         string
		allowedNamespaces []string
		expected          bool
	}{
		{"default", []string{""}, false},
		{"default", []string{"default"}, true},
		{"ns-123", []string{"foobar", "ns-*"}, true},
		{"ns", []string{"foobar", "ns-*"}, false},
		{"random-namespace", []string{"*"}, true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v/%v", test.namespace, test.allowedNamespaces), func(t *testing.T) {
			allowed := checkNamespace(test.namespace, test.allowedNamespaces)
			if allowed != test.expected {
				t.Errorf("got %v, expected %v", allowed, test.expected)
			}
		})
	}
}
