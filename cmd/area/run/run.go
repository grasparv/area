package run

import (
	"fmt"

	"github.com/grasparv/area/v2/cmd/area/userun"
	"github.com/grasparv/area/v2/internal/area"
)

func Run(n *area.Area, conf *area.Config) error {
	container, err := userun.Run(n, conf)
	if err != nil {
		return err
	}

	const guestBinary = "/bin/area-guest"
	exists, err := container.FileExists(guestBinary)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("%s does not exist in %s", guestBinary, n.ContainerName)
	}

	return container.ExecInteractive([]string{guestBinary})
}
