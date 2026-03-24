package userun

import (
	"os"
	"path/filepath"

	"github.com/grasparv/area/v2/pkg/podman"
)

func addHostResourceDisplay(spec *podman.ContainerSpec) {
	display := os.Getenv("DISPLAY")
	if display != "" {
		spec.Envs = append(spec.Envs, "DISPLAY="+display)
		if _, err := os.Stat("/tmp/.X11-unix"); err == nil {
			spec.Binds = append(spec.Binds, "/tmp/.X11-unix:/tmp/.X11-unix:rw")
		}
	}
}

func addHostResourceDRI(spec *podman.ContainerSpec) {
	if _, err := os.Stat("/dev/dri"); err == nil {
		spec.Binds = append(spec.Binds, "/dev/dri:/dev/dri")
	}
}

func addHostResourceXdg(spec *podman.ContainerSpec) bool {
	xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if xdgRuntimeDir == "" {
		return false
	}

	spec.Envs = append(spec.Envs, "XDG_RUNTIME_DIR="+xdgRuntimeDir)
	return true
}

func addHostResourceXdgWayland(spec *podman.ContainerSpec) {
	xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR")
	waylandDisplay := os.Getenv("WAYLAND_DISPLAY")
	if waylandDisplay != "" {
		spec.Envs = append(spec.Envs, "WAYLAND_DISPLAY="+waylandDisplay)
		waylandSocketPath := filepath.Join(xdgRuntimeDir, waylandDisplay)
		if _, err := os.Stat(waylandSocketPath); err == nil {
			spec.Binds = append(spec.Binds, waylandSocketPath+":"+waylandSocketPath)
		}
	}
}

func addHostResourceXdgPulse(spec *podman.ContainerSpec) {
	xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR")
	pulseSocketPath := filepath.Join(xdgRuntimeDir, "pulse/native")
	if _, err := os.Stat(pulseSocketPath); err == nil {
		spec.Binds = append(spec.Binds, pulseSocketPath+":"+pulseSocketPath)
		spec.Envs = append(spec.Envs, "PULSE_SERVER=unix:"+pulseSocketPath)
	}
}

func addHostResourceXdgPipewire(spec *podman.ContainerSpec) {
	xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR")
	pipewireSocketPath := filepath.Join(xdgRuntimeDir, "pipewire-0")
	if _, err := os.Stat(pipewireSocketPath); err == nil {
		spec.Binds = append(spec.Binds, pipewireSocketPath+":"+pipewireSocketPath)
	}
}
