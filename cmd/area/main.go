package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/grasparv/area/v2/cmd/area/ls"
	"github.com/grasparv/area/v2/cmd/area/rm"
	"github.com/grasparv/area/v2/cmd/area/run"
	"github.com/grasparv/area/v2/cmd/area/use"
	"github.com/grasparv/area/v2/internal/area"
	"github.com/grasparv/xflag/v2"
)

type UseCmd struct {
	XFlag string  `xflag:"use|Create or enter an area"`
	Name  string  `xflag:"Area name"`
	Image *string `xflag:"docker.io/library/debian:stable-slim|Image name"`
}

type RunCmd struct {
	XFlag string  `xflag:"run|Run /bin/area-guest in an area"`
	Name  string  `xflag:"Area name"`
	Image *string `xflag:"docker.io/library/debian:stable-slim|Image name"`
}

type ListCmd struct {
	XFlag string `xflag:"ls|List all areas"`
}

type RmCmd struct {
	XFlag string `xflag:"rm|Remove an area"`
	Name  string `xflag:"Area name"`
}

func main() {
	if _, err := exec.LookPath("podman"); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	commands := []any{
		UseCmd{},
		RunCmd{},
		ListCmd{},
		RmCmd{},
	}

	command, err := xflag.Parse(commands, os.Args)
	if err != nil {
		fmt.Print(err)
		return
	}

	conf := area.MakeConfig()

	switch c := command.(type) {
	case *UseCmd:
		n, err := area.New(c.Name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return
		}
		conf.ContainerImage = c.Image
		err = use.Run(n, conf)
	case *RunCmd:
		n, err := area.New(c.Name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return
		}
		conf.ContainerImage = c.Image
		err = run.Run(n, conf)
	case *ListCmd:
		err = ls.Run(nil, conf)
	case *RmCmd:
		n, err := area.New(c.Name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return
		}
		err = rm.Run(n, conf)
	default:
		fmt.Print(xflag.GetUsage(commands))
		return
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
