package rm

import (
	"fmt"

	"github.com/grasparv/area/v2/internal/area"
	"github.com/grasparv/area/v2/pkg/podman"
)

func Run(n *area.Area, conf *area.Config) error {
	p := podman.New()

	container := p.NewContainer(n.ContainerName)
	exists, err := container.Exists()
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("area %q does not exist", n.Name)
	}

	return container.Remove()
}
