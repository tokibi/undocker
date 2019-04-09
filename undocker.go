package undocker

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/tokibi/undocker/internal/untar"
)

type Undocker struct {
	Out, Err io.Writer
}

type Options struct {
	RegistryURL  string
	RegistryUser string
	RegistryPass string
	Extract      untar.Options
}

func (u Undocker) Extract(repo, tag, dest string, opts Options) error {
	source, err := createSource(opts)
	if err != nil {
		return err
	}
	if err := source.Image(repo, tag).Extract(dest, opts.Extract); err != nil {
		return err
	}
	return nil
}

func (u Undocker) Config(repo, tag string, opts Options) error {
	source, err := createSource(opts)
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

func createSource(opts Options) (src Source, err error) {
	url := opts.RegistryURL
	user := opts.RegistryUser
	pass := opts.RegistryPass

	if url != "" {
		src, err = NewRegistry(url, user, pass)
	} else {
		src, err = NewDockerAPI()
	}
	return
}
