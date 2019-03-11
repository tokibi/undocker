package undocker

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/docker/distribution"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

type Registry struct {
	URL      string
	Username string
	Password string
	session  *registry.Registry
}

func NewRegistry(url, username, password string) (*Registry, error) {
	r := &Registry{
		URL:      url,
		Username: username,
		Password: password,
	}
	if err := r.Authorize(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Registry) Authorize() error {
	sess, err := registry.New(r.URL, r.Username, r.Password)
	if err != nil {
		return err
	}
	r.session = sess
	return nil
}

func (r Registry) Exists(repository, tag string) bool {
	err := r.Find(repository, tag)
	if err != nil {
		return false
	}
	return true
}

func (r Registry) Find(repository, tag string) error {
	tags, err := r.session.Tags(repository)
	if err != nil {
		return errors.Wrap(err, "Repository not found")
	}
	for _, t := range tags {
		if t == tag {
			return nil
		}
	}
	return errors.New("Tag not found")
}

func (r Registry) LayerBlobs(repository, tag string) ([]io.Reader, error) {
	blobs := []io.Reader{}
	layers, err := r.Layers(repository, tag)
	if err != nil {
		return nil, err
	}
	for _, layer := range layers {
		blob, err := r.Blob(repository, layer.Digest)
		if err != nil {
			return nil, err
		}
		blobs = append(blobs, blob)
	}
	return blobs, nil
}

func (r Registry) Layers(repository, tag string) ([]distribution.Descriptor, error) {
	manifest, err := r.session.ManifestV2(repository, tag)
	if err != nil {
		return nil, err
	}
	return manifest.Layers, nil
}

func (r Registry) Blob(repository string, digest digest.Digest) (io.Reader, error) {
	blob, err := r.session.DownloadBlob(repository, digest)
	if err != nil {
		return nil, err
	}

	// Blob fetched from Docker Registry is compressed with gzip.
	gr, err := gzip.NewReader(blob)
	if err != nil {
		return nil, err
	}
	return gr, nil
}

func (r Registry) Image(repository, tag string) Image {
	return Image{
		Source:     r,
		Repository: repository,
		Tag:        tag,
	}
}
