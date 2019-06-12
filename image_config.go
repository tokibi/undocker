package undocker

import "time"

type ImageConfig struct {
	Architecture    string          `json:"architecture"`
	Config          Config          `json:"config"`
	Container       string          `json:"container"`
	ContainerConfig ContainerConfig `json:"container_config"`
	Created         time.Time       `json:"created"`
	DockerVersion   string          `json:"docker_version"`
	History         History         `json:"history"`
	OS              string          `json:"os"`
	Rootfs          Rootfs          `json:"rootfs"`
}

type Config struct {
	Hostname     string            `json:"Hostname"`
	Domainname   string            `json:"Domainname"`
	User         string            `json:"User"`
	AttachStdin  bool              `json:"AttachStdin"`
	AttachStdout bool              `json:"AttachStdout"`
	AttachStderr bool              `json:"AttachStderr"`
	Tty          bool              `json:"Tty"`
	OpenStdin    bool              `json:"OpenStdin"`
	StdinOnce    bool              `json:"StdinOnce"`
	Env          []string          `json:"Env"`
	Cmd          []string          `json:"Cmd"`
	ArgsEscaped  bool              `json:"ArgsEscaped"`
	Image        string            `json:"Image"`
	Volumes      []string          `json:"Volumes"`
	WorkingDir   string            `json:"WorkingDir"`
	Entrypoint   []string          `json:"Entrypoint"`
	OnBuild      []string          `json:"OnBuild"`
	Labels       map[string]string `json:"Labels"`
}

type ContainerConfig struct {
	Hostname     string            `json:"Hostname"`
	Domainname   string            `json:"Domainname"`
	User         string            `json:"User"`
	AttachStdin  bool              `json:"AttachStdin"`
	AttachStdout bool              `json:"AttachStdout"`
	AttachStderr bool              `json:"AttachStderr"`
	Tty          bool              `json:"Tty"`
	OpenStdin    bool              `json:"OpenStdin"`
	StdinOnce    bool              `json:"StdinOnce"`
	Env          []string          `json:"Env"`
	Cmd          []string          `json:"Cmd"`
	ArgsEscaped  bool              `json:"ArgsEscaped"`
	Image        string            `json:"Image"`
	Volumes      []string          `json:"Volumes"`
	WorkingDir   string            `json:"WorkingDir"`
	Entrypoint   []string          `json:"Entrypoint"`
	OnBuild      []string          `json:"OnBuild"`
	Labels       map[string]string `json:"Labels"`
}

type History []struct {
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	EmptyLayer bool      `json:"empty_layer,omitempty"`
}

type Rootfs struct {
	Type    string   `json:"type"`
	DiffIds []string `json:"diff_ids"`
}
