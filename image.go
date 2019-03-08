package undocker

import (
	"io"

	"github.com/pkg/errors"

	"github.com/docker/distribution"
	"github.com/opencontainers/go-digest"
	"github.com/tokibi/undocker/internal/untar"
)

type Source interface {
	Find(repository, tag string) error
	Layers(repository, tag string) ([]distribution.Descriptor, error)
	Blob(repository string, digest digest.Digest) (io.ReadCloser, error)
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
	layers, err := i.Layers()
	if err != nil {
		return err
	}
	for _, layer := range layers {
		reader, err := i.Blob(layer.Digest)
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

func (i Image) Exists() bool {
	if err := i.Source.Find(i.Repository, i.Tag); err != nil {
		return false
	}
	return true
}

func (i Image) Layers() ([]distribution.Descriptor, error) {
	return i.Source.Layers(i.Repository, i.Tag)
}

func (i Image) Blob(digest digest.Digest) (io.ReadCloser, error) {
	reader, err := i.Source.Blob(i.Repository, digest)
	if err != nil {
		return nil, err
	}
	return reader, nil
}
