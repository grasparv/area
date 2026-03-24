package userun

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/grasparv/area/v2/internal/area"
	"github.com/grasparv/area/v2/pkg/podman"
)

func Run(n *area.Area, conf *area.Config) (*podman.Container, error) {
	ctrl := podman.New()
	container := ctrl.NewContainer(n.ContainerName)

	exists, err := container.Exists()
	if err != nil {
		return nil, err
	}
	if !exists {
		container, err = create(ctrl, n, conf)
		if err != nil {
			return nil, err
		}
	}

	running, err := container.Running()
	if err != nil {
		return nil, err
	}
	if !running {
		if err := container.Start(); err != nil {
			return nil, err
		}
	}

	areaBinaryPath, err := getAreaBinaryPath()
	if err != nil {
		return nil, err
	}

	if err := container.CopyFromHost(areaBinaryPath, conf.GuestAreaBinaryPath); err != nil {
		return nil, err
	}

	if err := container.Chmod(conf.GuestAreaBinaryPath, "+x"); err != nil {
		return nil, err
	}

	return container, nil
}

func create(ctrl *podman.PodmanController, n *area.Area, conf *area.Config) (*podman.Container, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	homeDir, err = filepath.Abs(homeDir)
	if err != nil {
		return nil, err
	}

	mountForUser, err := getMountNameForUser()
	if err != nil {
		return nil, err
	}

	spec := podman.ContainerSpec{
		Name:     n.ContainerName,
		Hostname: n.Name,
		Image:    *conf.ContainerImage,
		Labels: map[string]string{
			conf.ManagedLabelKey: "true",
			conf.NameLabelKey:    n.Name,
		},
		SecurityOpts: []string{"label=disable"},
		Network:      "host",
		Binds: []string{
			"/:/mnt/host:rslave",
			homeDir + ":/home/" + mountForUser + ":rslave",
		},
		Command: []string{"sleep", "infinity"},
	}

	// Here we fix podman-in-podman socket paths.
	// The tricky thing is that inside the guest, we need to set the host's
	// socket path.
	if !runningInsideGuest() {
		if err := startPodmanUserSocket(); err != nil {
			return nil, err
		}
		pathPodmanSocketOnHost, err := getPodmanSocketHost()
		if err != nil {
			return nil, err
		}
		_, err = os.Stat(pathPodmanSocketOnHost)
		if err != nil {
			return nil, err
		}
		// We are on host.
		// -v /run/user/1000/podman/podman.sock:/run/podman/podman.sock
		// CONTAINER_HOST=unix:///run/podman/podman.sock
		// AREA_HOST_SOCKET_PATH=/run/user/1000/podman/podman.sock
		spec.Binds = append(spec.Binds, pathPodmanSocketOnHost+":"+pathPodmanSocketGuest)
		spec.Envs = append(spec.Envs, "CONTAINER_HOST=unix://"+pathPodmanSocketGuest)
		spec.Envs = append(spec.Envs, "AREA_HOST_SOCKET_PATH="+pathPodmanSocketOnHost)
	} else {
		// We are on guest.

		// -v /run/user/1000/podman/podman.sock:/run/podman/podman.sock
		// CONTAINER_HOST=unix:///run/podman/podman.sock
		// AREA_HOST_SOCKET_PATH=/run/user/1000/podman/podman.sock
		//
		// i.e. the guest needs to use the host's socket path that we get
		// from the AREA_HOST_SOCKET_PATH variable
		pathPodmanSocketOnHost := os.Getenv("AREA_HOST_SOCKET_PATH")
		if pathPodmanSocketOnHost == "" {
			return nil, errors.New("unknown host podman socket path")
		}

		spec.Binds = append(spec.Binds, pathPodmanSocketOnHost+":"+pathPodmanSocketGuest)
		spec.Envs = append(spec.Envs, "CONTAINER_HOST=unix://"+pathPodmanSocketGuest)
		spec.Envs = append(spec.Envs, "AREA_HOST_SOCKET_PATH="+pathPodmanSocketOnHost)
	}

	addHostResourceDisplay(&spec)
	addHostResourceDRI(&spec)
	if addHostResourceXdg(&spec) {
		addHostResourceXdgWayland(&spec)
		addHostResourceXdgPulse(&spec)
		addHostResourceXdgPipewire(&spec)
	}

	container := ctrl.NewContainer(n.ContainerName)
	return container, container.Create(spec)
}

func runningInsideGuest() bool {
	_, err := os.Stat(pathPodmanSocketGuest)
	return err == nil
}
