package area

import (
	"fmt"
	"strings"
)

const (
	containerPrefix     = "area-"
	guestAreaBinaryPath = "/bin/area"
	nameLabelKey        = "dev.area.name"
	managedLabelKey     = "dev.area.managed"
)

type Area struct {
	Name          string
	ContainerName string
}

type Config struct {
	ContainerPrefix     string
	ContainerImage      *string
	GuestAreaBinaryPath string
	NameLabelKey        string
	ManagedLabelKey     string
}

func MakeConfig() *Config {
	return &Config{
		ContainerPrefix:     containerPrefix,
		GuestAreaBinaryPath: guestAreaBinaryPath,
		NameLabelKey:        nameLabelKey,
		ManagedLabelKey:     managedLabelKey,
	}
}

func New(name string) (*Area, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	for _, r := range name {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= 'A' && r <= 'Z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		return nil, fmt.Errorf("invalid name %q: only letters and numbers are allowed", name)
	}

	return &Area{
		Name:          name,
		ContainerName: containerPrefix + name,
	}, nil
}
