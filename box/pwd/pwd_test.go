package pwd

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/zooyer/gobox/box"
)

// Helper function to call system `pwd` with options
func systemPwd(args ...string) (string, error) {
	cmd := exec.Command("pwd", args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}

func TestPwd(t *testing.T) {
	tests := []struct {
		Name         string
		Args         []string
		ExpectedCode int
	}{
		{"DefaultPhysical", []string{}, 0},
		{"Logical", []string{"-L"}, 0},
		{"Physical", []string{"-P"}, 0},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// Get system `pwd` output
			expectedOutput, err := systemPwd(test.Args...)
			if err != nil {
				t.Fatalf("Failed to run system pwd: %v", err)
			}

			// Run our Pwd implementation
			var (
				stdout, stderr bytes.Buffer
				option         = box.Option{
					Args:   append([]string{"echo"}, test.Args...),
					Stdout: &stdout,
					Stderr: &stderr,
				}
			)

			code := Pwd(context.Background(), option)
			if code != test.ExpectedCode {
				t.Errorf("[%s] Unexpected exit code: got %d, want %d", test.Name, code, test.ExpectedCode)
			}

			// Compare output
			output := strings.TrimSpace(stdout.String())
			if output != expectedOutput {
				t.Errorf("[%s] Mismatch:\nExpected: %s\nGot: %s", test.Name, expectedOutput, output)
			}
		})
	}
}
