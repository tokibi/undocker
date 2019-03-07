package undocker

import (
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

func (r *Registry) Authorize() (*registry.Registry, error) {
	if r.session != nil {
		return r.session, nil
	}

	reg, err := registry.New(r.URL, r.Username, r.Password)
	if err != nil {
		return nil, err
	}
	r.session = reg
	return reg, nil
}

func (r Registry) Find(repo, tag string) error {
	sess, err := r.Authorize()
	if err != nil {
		return err
	}
	tags, err := sess.Tags(repo)
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

func (r Registry) Layers(repo, tag string) ([]distribution.Descriptor, error) {
	sess, err := r.Authorize()
	if err != nil {
		return nil, err
	}
	manifest, err := sess.ManifestV2(repo, tag)
	if err != nil {
		return nil, err
	}
	return manifest.Layers, nil
}

func (r Registry) Blob(repo string, digest digest.Digest) (io.ReadCloser, error) {
	sess, err := r.Authorize()
	if err != nil {
		return nil, err
	}
	return sess.DownloadBlob(repo, digest)
}

func (r Registry) Image(repo, tag string) Image {
	return Image{
		Source:     r,
		Repository: repo,
		Tag:        tag,
	}
}
