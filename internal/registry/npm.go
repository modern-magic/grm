package registry

import (
	"os"
	"os/exec"
)

func WriteNpm(uri string) error {
	args := []string{"config", "set", "registry", uri}
	cmd := exec.Command("npm", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ReadNpm() (string, error) {
	args := []string{"config", "get", "registry"}
	cmd := exec.Command("npm", args...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
