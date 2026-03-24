package userun

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func getMountNameForUser() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	name := strings.ToLower(strings.TrimSpace(currentUser.Username))

	normalized := make([]rune, 0, len(name))
	for _, r := range name {
		if r >= 'a' && r <= 'z' {
			normalized = append(normalized, r)
			continue
		}
		if r >= '0' && r <= '9' {
			normalized = append(normalized, r)
		}
	}

	if len(normalized) == 0 {
		return "", fmt.Errorf("user mount name cannot be empty")
	}

	return string(normalized), nil
}

func getAreaBinaryPath() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	areaBinaryPath, err := filepath.EvalSymlinks(executablePath)
	if err != nil {
		return filepath.Clean(executablePath), nil
	}

	return areaBinaryPath, nil
}
