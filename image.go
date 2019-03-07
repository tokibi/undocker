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

func (i Image) Unpack(dir string) error {
	layers, err := i.Layers()
	if err != nil {
		return err
	}

	for _, layer := range layers {
		reader, err := i.Blob(layer.Descriptor().Digest)
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

func (i Image) Exists() error {
	return i.Source.Find(i.Repository, i.Tag)
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
