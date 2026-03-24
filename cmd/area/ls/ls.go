package ls

import (
	"fmt"
	"strings"

	"github.com/grasparv/area/v2/internal/area"
	"github.com/grasparv/area/v2/pkg/podman"
)

func Run(n *area.Area, conf *area.Config) error {
	p := podman.New()

	names, err := p.ContainerNames(conf.ManagedLabelKey)
	if err != nil {
		return err
	}

	for _, name := range names {
		fmt.Println(strings.TrimPrefix(name, conf.ContainerPrefix))
	}

	return nil
}
