package undocker

import (
	"io"
	"strings"

	"github.com/urfave/cli"
)

type Undocker struct {
	Out, Err io.Writer
}

func (u Undocker) Extract(c *cli.Context) error {
	url := c.GlobalString("registry-url")
	user := c.GlobalString("registry-user")
	pass := c.GlobalString("registry-pass")

	var source Source
	var err error
	if url != "" {
		source, err = NewRegistry(url, user, pass)
	} else {
		source, err = NewDockerAPI()
	}
	if err != nil {
		return err
	}

	ref := strings.SplitN(c.Args().Get(0), ":", 2)
	dst := c.Args().Get(1)
	if dst == "" {
		dst = "."
	}
	if err := source.Image(ref[0], ref[1]).Unpack(dst); err != nil {
		return err
	}
	return nil
}
