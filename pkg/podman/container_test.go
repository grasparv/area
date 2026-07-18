package podman

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestContainerCreateAddsDevices(t *testing.T) {
	tmpDir := t.TempDir()
	argsPath := filepath.Join(tmpDir, "args")
	podmanPath := filepath.Join(tmpDir, "podman")
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$PODMAN_TEST_ARGS\"\n"
	if err := os.WriteFile(podmanPath, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	t.Setenv("PATH", tmpDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("PODMAN_TEST_ARGS", argsPath)

	container := New().NewContainer("test")
	err := container.Create(ContainerSpec{
		Name:     "test",
		Hostname: "test",
		Image:    "example.invalid/image",
		Devices:  []string{"nvidia.com/gpu=all"},
	})
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(argsPath)
	if err != nil {
		t.Fatal(err)
	}
	args := strings.Fields(string(data))
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--device nvidia.com/gpu=all") {
		t.Fatalf("podman arguments do not contain NVIDIA CDI device: %q", joined)
	}
}
