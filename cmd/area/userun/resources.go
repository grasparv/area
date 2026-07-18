package userun

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/grasparv/area/v2/pkg/podman"
)

const nvidiaCDIDevice = "nvidia.com/gpu=all"

var (
	nvidiaCDIKindPattern = regexp.MustCompile(`(?m)^\s*(?:"kind"|kind)\s*:\s*["']?nvidia\.com/gpu["']?\s*,?\s*$`)
	nvidiaCDIAllPattern  = regexp.MustCompile(`(?m)^\s*(?:-\s*)?(?:"name"|name)\s*:\s*["']?all["']?\s*,?\s*$`)
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

func addHostResourceNVIDIA(spec *podman.ContainerSpec) {
	addHostResourceNVIDIAAtRoot(spec, hostFilesystemRoot())
}

func addHostResourceNVIDIAAtRoot(spec *podman.ContainerSpec, root string) {
	if !hostHasNVIDIAGPU(root) || !hostHasNVIDIACDI(root) {
		return
	}

	spec.Devices = append(spec.Devices, nvidiaCDIDevice)
}

func hostFilesystemRoot() string {
	if runningInsideGuest() {
		return "/mnt/host"
	}
	return "/"
}

func hostHasNVIDIAGPU(root string) bool {
	if _, err := os.Stat(filepath.Join(root, "dev", "nvidiactl")); err != nil {
		return false
	}

	matches, err := filepath.Glob(filepath.Join(root, "dev", "nvidia[0-9]*"))
	if err != nil {
		return false
	}

	for _, match := range matches {
		index := strings.TrimPrefix(filepath.Base(match), "nvidia")
		if index != "" && strings.IndexFunc(index, func(r rune) bool {
			return r < '0' || r > '9'
		}) == -1 {
			return true
		}
	}

	return false
}

func hostHasNVIDIACDI(root string) bool {
	for _, dir := range []string{
		filepath.Join(root, "etc", "cdi"),
		filepath.Join(root, "var", "run", "cdi"),
	} {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() || !isCDISpecFilename(entry.Name()) {
				continue
			}

			data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			if err != nil {
				continue
			}

			if nvidiaCDIKindPattern.Match(data) && nvidiaCDIAllPattern.Match(data) {
				return true
			}
		}
	}

	return false
}

func isCDISpecFilename(name string) bool {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".json", ".yaml", ".yml":
		return true
	default:
		return false
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
