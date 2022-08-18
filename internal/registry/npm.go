package registry

import (
	"os"
	"os/exec"
)

func WriteNpm(uri string) error {
	args := []string{"config", "set", "registry", uri}
	cmd := exec.Command("npm", args...)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func ReadNpm() string {
	args := []string{"config", "get", "registry"}
	cmd := exec.Command("npm", args...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return err.Error()
	}
	return string(out)
}
