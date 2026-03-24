package podman

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type PodmanController struct {
}

func New() *PodmanController {
	return &PodmanController{}
}

func (p *PodmanController) NewContainer(name string) *Container {
	return &Container{podman: p, name: name}
}

func (p *PodmanController) ContainerNames(managedLabelKey string) ([]string, error) {
	managedLabelFilter := "label=" + managedLabelKey + "=true"
	args := []string{
		"--all",
		"--filter", managedLabelFilter,
		"--format", "{{.Names}}",
	}

	output, err := p.runCapture(append([]string{"ps"}, args...)...)
	if err != nil {
		return nil, err
	}

	output = strings.TrimSpace(output)
	if output == "" {
		return nil, nil
	}

	rawNames := strings.Split(output, "\n")
	names := make([]string, 0, len(rawNames))
	for _, rawName := range rawNames {
		name := strings.TrimSpace(rawName)
		if name != "" {
			names = append(names, name)
		}
	}

	return names, nil
}

func (p *PodmanController) run(args ...string) error {
	cmd := p.newCommand(args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *PodmanController) runCapture(args ...string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := p.newCommand(args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("podman %s failed: %s", strings.Join(args, " "), strings.TrimSpace(stderr.String()))
		}
		return "", err
	}

	return stdout.String(), nil
}

func (p *PodmanController) newCommand(args ...string) *exec.Cmd {
	cmd := exec.Command("podman", args...)
	cmd.Env = os.Environ()
	return cmd
}
