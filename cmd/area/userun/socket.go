package userun

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	pathPodmanSocketGuest = "/run/podman/podman.sock"
)

func startPodmanUserSocket() error {
	if _, err := exec.LookPath("systemctl"); err != nil {
		return fmt.Errorf("podman socket not available and systemctl not found")
	}

	fmt.Println("Starting podman.socket (required for area)...")

	cmd := exec.Command("systemctl", "--user", "enable", "--now", "podman.socket")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to enable podman.socket: %w", err)
	}

	socketPath := fmt.Sprintf("/run/user/%d/podman/podman.sock", os.Getuid())
	return waitForPodmanSocket(socketPath, 5*time.Second)
}

func waitForPodmanSocket(socketPath string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	tr := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", socketPath)
		},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   2 * time.Second,
	}

	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("podman socket not ready after %s", timeout)
		}

		resp, err := client.Get("http://d/_ping")
		if err == nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			if resp.StatusCode == 200 {
				return nil
			}
		}

		time.Sleep(200 * time.Millisecond)
	}
}

func getPodmanSocketHost() (string, error) {
	if xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR"); xdgRuntimeDir != "" {
		socketPath := filepath.Join(xdgRuntimeDir, "podman", "podman.sock")
		if _, err := os.Stat(socketPath); err == nil {
			return socketPath, nil
		}
	}

	uid := os.Getuid()
	fallback := filepath.Join("/run/user", fmt.Sprintf("%d", uid), "podman", "podman.sock")
	if _, err := os.Stat(fallback); err == nil {
		return fallback, nil
	}

	return "", errors.New("could not determine host podman socket path")
}
