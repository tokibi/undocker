package undocker

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"

	"github.com/urfave/cli"
)

type Undocker struct {
	Out, Err io.Writer
}

func (u Undocker) Extract(c *cli.Context) error {
	source, err := createSource(c)
	if err != nil {
		return err
	}
	repo, tag, err := parseReference(c.Args().Get(0))
	if err != nil {
		return err
	}
	dst := c.Args().Get(1)
	if dst == "" {
		dst = "."
	}

	if err := source.Image(repo, tag).Unpack(dst); err != nil {
		return err
	}
	return nil
}

func (u Undocker) Config(c *cli.Context) error {
	source, err := createSource(c)
	if err != nil {
		return err
	}
	repo, tag, err := parseReference(c.Args().Get(0))
	if err != nil {
		return err
	}

	config, err := source.Image(repo, tag).Config()
	if err != nil {
		return err
	}
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	fmt.Fprintln(u.Out, string(data))
	return nil
}

func createSource(c *cli.Context) (src Source, err error) {
	url := c.GlobalString("registry-url")
	user := c.GlobalString("registry-user")
	pass := c.GlobalString("registry-pass")

	if url != "" {
		src, err = NewRegistry(url, user, pass)
	} else {
		src, err = NewDockerAPI()
	}
	return
}

func parseReference(arg string) (repository, tag string, err error) {
	ref := strings.SplitN(arg, ":", 2)
	if ref[0] == "" {
		return "", "", errors.New("Invalid image")
	}
	if len(ref) < 2 {
		return ref[0], "latest", nil
	}
	return ref[0], ref[1], nil
}
