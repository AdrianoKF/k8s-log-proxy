package logs

import (
	"testing"
)

func TestParseUrl(t *testing.T) {
	tests := []struct {
		path        string
		expected    LogTarget
		expectedErr bool
	}{
		{"/logs/default/pod", LogTarget{Namespace: "default", Pod: "pod"}, false},
		{"/logs/default/pod/container", LogTarget{Namespace: "default", Pod: "pod", Container: "container"}, false},
		{"/logs//pod", LogTarget{}, true},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			actual, err := ParseUrl(test.path)
			if test.expectedErr {
				if err == nil {
					t.Error("expected error but got result", actual)
				}
			} else {
				if err != nil {
					t.Error("got unexpected error:", err)
				} else if actual != test.expected {
					t.Errorf("actual != expected: %+v != %+v", actual, test.expected)
				}
			}
		})
	}
}
