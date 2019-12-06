package undocker

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/moby/moby/client"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type DockerAPI struct {
	client  *client.Client
	context context.Context
}

func NewDockerAPI() (*DockerAPI, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &DockerAPI{
		client:  cli,
		context: context.Background(),
	}, nil
}

func (api DockerAPI) Exists(repository, tag string) bool {
	_, err := api.Find(repository, tag)
	if err != nil {
		return false
	}
	return true
}

func (api DockerAPI) Find(repository, tag string) (string, error) {
	args := filters.NewArgs()
	args.Add("reference", fmt.Sprintf("%s:%s", repository, tag))

	list, err := api.client.ImageList(api.context, types.ImageListOptions{
		Filters: args,
	})
	if err != nil {
		return "", err
	}
	if len(list) == 0 {
		return "", errors.New("Image not found")
	}
	return list[0].ID, nil
}

func (api DockerAPI) LayerBlobs(repository, tag string) ([]io.Reader, error) {
	blob, err := api.ImageBlob(repository, tag)
	if err != nil {
		return nil, err
	}
	return blob.LayerBlobs()
}

func (api DockerAPI) Config(repository, tag string) ([]byte, error) {
	blob, err := api.ImageBlob(repository, tag)
	if err != nil {
		return nil, err
	}
	return blob.Config()
}

func (api DockerAPI) Image(repository, tag string) Image {
	return Image{
		Source:     api,
		Repository: repository,
		Tag:        tag,
	}
}

func (api DockerAPI) ImageBlob(repository, tag string) (*ImageBlob, error) {
	id, err := api.Find(repository, tag)
	if err != nil {
		return nil, err
	}
	blob, err := api.client.ImageSave(api.context, []string{id})
	if err != nil {
		return nil, err
	}
	return &ImageBlob{Blob: blob}, err
}

func (api DockerAPI) CleanUp() error {
	return nil
}

type ImageBlob struct {
	Blob io.ReadCloser
}

func (i *ImageBlob) Manifest() (Manifest, error) {
	var manifest Manifest
	tr := tar.NewReader(i.Blob)
	for {
		f, err := tr.Next()
		if err == io.EOF {
			return manifest, errors.New("Manifest file not found")
		}
		if f.Name == "manifest.json" {
			break
		}
	}
	manifest, err := unmarshalManifest(tr)
	if err != nil {
		return manifest, err
	}
	return manifest, nil
}

func (i *ImageBlob) LayerBlobs() ([]io.Reader, error) {
	var manifest Manifest
	var blobs []io.Reader
	bufs := map[string]io.Reader{}

	tr := tar.NewReader(i.Blob)
	for {
		f, err := tr.Next()
		if err == io.EOF {
			break
		}
		if f.Name == "manifest.json" {
			manifest, err = unmarshalManifest(tr)
			if err != nil {
				return nil, err
			}
		}
		if strings.HasSuffix(f.Name, "/layer.tar") {
			buf := new(bytes.Buffer)
			io.Copy(buf, tr)
			bufs[f.Name] = buf
		}
	}
	// Make it the same order as Manifest.
	for _, layer := range manifest.Layers {
		blobs = append(blobs, bufs[layer])
	}
	return blobs, nil
}

func (i *ImageBlob) Config() ([]byte, error) {
	var manifest Manifest
	bufs := map[string]io.Reader{}

	tr := tar.NewReader(i.Blob)
	for {
		f, err := tr.Next()
		if err == io.EOF {
			break
		}
		if f.Name == "manifest.json" {
			manifest, err = unmarshalManifest(tr)
			if err != nil {
				return nil, err
			}
		}
		if strings.HasSuffix(f.Name, ".json") {
			buf := new(bytes.Buffer)
			io.Copy(buf, tr)
			bufs[f.Name] = buf
		}
	}
	return ioutil.ReadAll(bufs[manifest.Config])
}

type Manifest struct {
	Config   string   `json:"Config"`
	RepoTags []string `json:"RepoTags"`
	Layers   []string `json:"Layers"`
}

func unmarshalManifest(r io.Reader) (Manifest, error) {
	var manifest Manifest

	buf := new(bytes.Buffer)
	io.Copy(buf, r)
	m := new([]Manifest)
	if err := json.Unmarshal(buf.Bytes(), m); err != nil {
		return manifest, err
	}
	manifest = (*m)[0]
	return manifest, nil
}
