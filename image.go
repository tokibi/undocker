package undocker

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"

	"github.com/tokibi/undocker/internal/untar"
)

type Source interface {
	Config(repository, tag string) ([]byte, error)
	Exists(repository, tag string) bool
	LayerBlobs(repository, tag string) ([]io.Reader, error)
	Image(repository, tag string) Image
}

type Image struct {
	Source     Source
	Repository string
	Tag        string
}

func (i Image) Unpack(dir string) error {
	if !i.Exists() {
		return errors.New("Image not found")
	}
	layerBlobs, err := i.LayerBlobs()
	if err != nil {
		return err
	}
	for _, blob := range layerBlobs {
		untar.Untar(blob, dir)
	}
	return nil
}

func (i Image) Config() (*ImageConfig, error) {
	bytes, err := i.Source.Config(i.Repository, i.Tag)
	if err != nil {
		return nil, err
	}
	config := new(ImageConfig)
	if err := json.Unmarshal(bytes, config); err != nil {
		return nil, err
	}
	return config, nil
}

func (i Image) Exists() bool {
	if i.Source.Exists(i.Repository, i.Tag) {
		return true
	}
	return false
}

func (i Image) LayerBlobs() ([]io.Reader, error) {
	return i.Source.LayerBlobs(i.Repository, i.Tag)
}

