package use

import (
	"github.com/grasparv/area/v2/cmd/area/userun"
	"github.com/grasparv/area/v2/internal/area"
)

func Run(n *area.Area, conf *area.Config) error {
	container, err := userun.Run(n, conf)
	if err != nil {
		return err
	}

	const guestBinary = "/bin/sh"
	return container.ExecInteractive([]string{guestBinary})
}
