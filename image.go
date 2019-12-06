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
	CleanUp() error
}

type Image struct {
	Source     Source
	Repository string
	Tag        string
}

// Extract extracts docker image as rootfs to the specified directory
func (i Image) Extract(dir string, overwriteSymlink bool) error {
	if !i.Exists() {
		return errors.New("Image not found")
	}
	layerBlobs, err := i.LayerBlobs()
	if err != nil {
		return err
	}
	for _, blob := range layerBlobs {
		err = untar.Untar(blob, dir, untar.Options{
			OverwriteSymlinkRefs: overwriteSymlink,
		})

		if err != nil {
			return err
		}
	}
	return nil
}

// Unpack is an alias for Extract()
func (i Image) Unpack(dir string, overwriteSymlink bool) error {
	return i.Extract(dir, overwriteSymlink)
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

// Exists check the images
func (i Image) Exists() bool {
	if i.Source.Exists(i.Repository, i.Tag) {
		return true
	}
	return false
}

// LayerBlobs return the layers of the image in order from the lower
func (i Image) LayerBlobs() ([]io.Reader, error) {
	return i.Source.LayerBlobs(i.Repository, i.Tag)
}
