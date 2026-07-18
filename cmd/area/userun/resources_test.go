package userun

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/grasparv/area/v2/pkg/podman"
)

func TestAddHostResourceNVIDIAAtRoot(t *testing.T) {
	tests := []struct {
		name       string
		withGPU    bool
		withCDI    bool
		cdiDir     string
		wantDevice bool
	}{
		{name: "neither available"},
		{name: "gpu only", withGPU: true},
		{name: "cdi only", withCDI: true, cdiDir: "etc/cdi"},
		{name: "etc cdi", withGPU: true, withCDI: true, cdiDir: "etc/cdi", wantDevice: true},
		{name: "runtime cdi", withGPU: true, withCDI: true, cdiDir: "var/run/cdi", wantDevice: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			if test.withGPU {
				createTestFile(t, filepath.Join(root, "dev", "nvidiactl"), "")
				createTestFile(t, filepath.Join(root, "dev", "nvidia0"), "")
			}
			if test.withCDI {
				createTestFile(t, filepath.Join(root, test.cdiDir, "nvidia.yaml"), `
cdiVersion: 0.5.0
kind: nvidia.com/gpu
devices:
  - name: all
`)
			}

			spec := &podman.ContainerSpec{}
			addHostResourceNVIDIAAtRoot(spec, root)

			if test.wantDevice {
				if len(spec.Devices) != 1 || spec.Devices[0] != nvidiaCDIDevice {
					t.Fatalf("devices = %#v, want [%q]", spec.Devices, nvidiaCDIDevice)
				}
			} else if len(spec.Devices) != 0 {
				t.Fatalf("devices = %#v, want none", spec.Devices)
			}
		})
	}
}

func TestHostHasNVIDIAGPUIgnoresNonNumericDeviceNames(t *testing.T) {
	root := t.TempDir()
	createTestFile(t, filepath.Join(root, "dev", "nvidiactl"), "")
	createTestFile(t, filepath.Join(root, "dev", "nvidia-modeset"), "")
	createTestFile(t, filepath.Join(root, "dev", "nvidia0bad"), "")

	if hostHasNVIDIAGPU(root) {
		t.Fatal("hostHasNVIDIAGPU returned true without a numeric GPU device")
	}
}

func TestHostHasNVIDIACDIIgnoresOtherKinds(t *testing.T) {
	root := t.TempDir()
	createTestFile(t, filepath.Join(root, "etc", "cdi", "other.yaml"), `
cdiVersion: 0.5.0
kind: example.com/device
devices:
  - name: all
`)

	if hostHasNVIDIACDI(root) {
		t.Fatal("hostHasNVIDIACDI returned true for a non-NVIDIA CDI spec")
	}
}

func createTestFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
