package undocker

import (
	"io"

	"github.com/docker/distribution"
	"github.com/opencontainers/go-digest"
	"github.com/tokibi/undocker/internal/untar"
)

type Source interface {
	Find(repo, tag string) error
	Layers(repo, tag string) ([]distribution.Descriptor, error)
	Blob(repo string, digest digest.Digest) (io.ReadCloser, error)
	Image(repo, tag string) Image
}

type Image struct {
	Source     Source
	Repository string
	Tag        string
}

func (image Image) Unpack(dir string) error {
	layers, err := image.Layers()
	if err != nil {
		return err
	}

	for _, layer := range layers {
		reader, err := image.Blob(layer.Descriptor().Digest)
		if err != nil {
			return err
		}
		if reader != nil {
			untar.Untar(reader, dir)
			reader.Close()
		}
	}

	return nil
}

func (image Image) Exists() error {
	return image.Source.Find(image.Repository, image.Tag)
}

func (image Image) Layers() ([]distribution.Descriptor, error) {
	return image.Source.Layers(image.Repository, image.Tag)
}

func (image Image) Blob(digest digest.Digest) (io.ReadCloser, error) {
	reader, err := image.Source.Blob(image.Repository, digest)
	if err != nil {
		return nil, err
	}
	return reader, nil
}
