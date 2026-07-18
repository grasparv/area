package podman

import (
	"errors"
	"os/exec"
	"sort"
	"strings"
)

type ContainerSpec struct {
	Name         string
	Hostname     string
	Image        string
	Labels       map[string]string
	SecurityOpts []string
	Devices      []string
	Network      string
	Binds        []string
	Envs         []string
	Command      []string
}

type Container struct {
	podman *PodmanController
	name   string
}

func (c *Container) Exists() (bool, error) {
	err := c.podman.run("container", "exists", c.name)
	if err == nil {
		return true, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return false, nil
	}

	return false, err
}

func (c *Container) Running() (bool, error) {
	output, err := c.podman.runCapture("inspect", "--format", "{{.State.Running}}", c.name)
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(output) == "true", nil
}

func (c *Container) FileExists(path string) (bool, error) {
	err := c.podman.run("exec", c.name, "test", "-e", path)
	if err == nil {
		return true, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return false, nil
	}

	return false, err
}

func (c *Container) Create(spec ContainerSpec) error {
	args := []string{"create", "--name", spec.Name, "--hostname", spec.Hostname}

	labelKeys := make([]string, 0, len(spec.Labels))
	for key := range spec.Labels {
		labelKeys = append(labelKeys, key)
	}
	sort.Strings(labelKeys)
	for _, key := range labelKeys {
		args = append(args, "--label", key+"="+spec.Labels[key])
	}
	for _, securityOpt := range spec.SecurityOpts {
		args = append(args, "--security-opt", securityOpt)
	}
	for _, device := range spec.Devices {
		args = append(args, "--device", device)
	}
	if spec.Network != "" {
		args = append(args, "--network", spec.Network)
	}
	for _, bind := range spec.Binds {
		args = append(args, "-v", bind)
	}
	for _, envVar := range spec.Envs {
		args = append(args, "-e", envVar)
	}

	args = append(args, spec.Image)
	args = append(args, spec.Command...)
	return c.podman.run(args...)
}

func (c *Container) Start() error {
	return c.podman.run("start", c.name)
}

func (c *Container) Remove() error {
	return c.podman.run("rm", "-f", c.name, "-t", "1")
}

func (c *Container) ExecInteractive(command []string) error {
	args := append([]string{"exec", "-it", c.name}, command...)
	return c.podman.run(args...)
}

func (c *Container) CopyFromHost(srcPath, dstPath string) error {
	return c.podman.run("cp", srcPath, c.name+":"+dstPath)
}

func (c *Container) Chmod(path, mode string) error {
	return c.podman.run("exec", "-u", "root", c.name, "chmod", mode, path)
}
