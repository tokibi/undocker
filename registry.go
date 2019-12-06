package undocker

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"time"
)

const dockerHubServiceName = "registry.docker.io"

type Registry struct {
	URL         *url.URL
	Username    string
	Password    string
	client      *registry.Registry
	isDockerHub bool
	tmpDir      string
}

func NewRegistry(baseURL, username, password, tmpRoot string) (*Registry, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	c, err := auth(u.String(), username, password)
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().Format("2006010215030405")
	tmpdir := filepath.Join(tmpRoot, timestamp)
	err = os.MkdirAll(tmpdir, 0666)
	if err != nil {
		return nil, err
	}

	return &Registry{
		URL:         u,
		Username:    username,
		Password:    password,
		client:      c,
		isDockerHub: isDockerHub(u),
		tmpDir:      tmpdir,
	}, nil
}

// auth returns client for docker registry
func auth(url, username, password string) (*registry.Registry, error) {
	c, err := registry.New(url, username, password)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r Registry) Exists(repository, tag string) bool {
	err := r.Find(repository, tag)
	if err != nil {
		return false
	}
	return true
}

func (r Registry) Manifest(repository, tag string) (*schema2.DeserializedManifest, error) {
	return r.client.ManifestV2(repository, tag)
}

func (r Registry) Find(repository, tag string) error {
	tags, err := r.client.Tags(repository)
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
		blob, err := r.ExtractedBlob(repository, layer.Digest)
		if err != nil {
			return nil, err
		}
		blobs = append(blobs, blob)
	}
	return blobs, nil
}

func (r Registry) Layers(repository, tag string) ([]distribution.Descriptor, error) {
	manifest, err := r.client.ManifestV2(repository, tag)
	if err != nil {
		return nil, err
	}
	return manifest.Layers, nil
}

func (r Registry) CleanUp() error {
	if err := os.RemoveAll(r.tmpDir); err != nil {
		return nil
	}
	return nil
}

func (r Registry) ExtractedBlob(repository string, digest digest.Digest) (io.Reader, error) {
	httpReader, err := r.client.DownloadBlob(repository, digest)
	if err != nil {
		return nil, err
	}

	if stat, err := os.Stat(r.tmpDir); err != nil {
		return nil, err
	} else if !stat.IsDir() {
		return nil, errors.Errorf("%s must be directory", r.tmpDir)
	}

	tmpFilePath := filepath.Join(r.tmpDir, digest.String())
	f, err := os.OpenFile(tmpFilePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(f, httpReader)
	if err != nil {
		return nil, err
	}

	err = httpReader.Close()
	if err != nil {
		return nil, err
	}

	err = httpReader.Close()
	if err != nil {
		return nil, err
	}

	blob, err := os.OpenFile(tmpFilePath, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}

	gr, err := gzip.NewReader(blob)
	if err != nil {
		return nil, err
	}

	return gr, nil
}

func (r Registry) Image(repository, tag string) Image {
	if r.isDockerHub {
		repository = complementOfficialRepoName(repository)
	}
	return Image{
		Source:     r,
		Repository: repository,
		Tag:        tag,
	}
}

func (r Registry) Config(repository, tag string) ([]byte, error) {
	manifest, err := r.Manifest(repository, tag)
	if err != nil {
		return nil, err
	}
	reader, err := r.client.DownloadBlob(repository, manifest.Config.Digest)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

func complementOfficialRepoName(repository string) string {
	if len(strings.Split(repository, "/")) == 1 {
		return "library/" + repository
	}
	return repository
}

func isDockerHub(url *url.URL) bool {
	u := *url
	u.Path = path.Join(u.Path, "v2")

	resp, err := http.Get(u.String() + "/")
	if err != nil {
		return false
	}

	h := resp.Header.Get("Www-Authenticate")
	if strings.Contains(h, fmt.Sprintf(`service="%s"`, dockerHubServiceName)) {
		return true
	}
	return false
}
