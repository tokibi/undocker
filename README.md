# undocker

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go library and command line tool for decomposing docker images.

## Command Use

```
NAME:
   undocker - Decompose docker images.

USAGE:
   undocker [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     extract, e  Extract to rootfs.
     config, c   Show image configuration.
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --registry-url value, -r value   docker registry url [$REGISTRY_URL]
   --registry-user value, -u value  docker registry login username [$REGISTRY_USER]
   --registry-pass value, -p value  docker registry login password [$REGISTRY_PASS]
   --help, -h                       show help
   --version, -v                    print the version
```

### Extract

Extract from local images.

```console
$ undocker extract busybox:latest ./image
$ ls ./image
bin/  dev/  etc/  home/  root/	tmp/  usr/  var/
```

Extract directly from docker registry.

```console
$ export REGISTRY_USER=xxx # optional
$ export REGISTRY_PASS=xxx # optional
$ undocker -r "https://registry-1.docker.io/" extract busybox:latest ./image
```

### Config

Show image config.

```console
$ undocker config busybox:latest | jq
{
  "architecture": "amd64",
  "config": {
    "Hostname": "",
    "Domainname": "",
    "User": "",
    "AttachStdin": false,
    "AttachStdout": false,
    "AttachStderr": false,
    "Tty": false,
    "OpenStdin": false,
    "StdinOnce": false,
    "Env": [
      "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
    ],
    "Cmd": [
      "sh"
    ],
...
```

## Library Use

### Extract

Extract from local images.

```go
func main() {
    dst := "./image"

    api, err := undocker.NewDockerAPI()
    if err != nil {
        log.Fatal(err)
    }
    api.Image("busybox", "latest").Unpack(dst)
}
```

Extract directly from docker registry.

```go
func main() {
    url := "https://registry-1.docker.io/"
    username := ""
    password := ""
    dst := "./image"

    registry, err := undocker.NewRegistry(url, username, password)
    if err != nil {
        log.Fatal(err)
    }
    registry.Image("busybox", "latest").Unpack(dst)
}
```

### Config

```go
func main() {
    api, _ := undocker.NewDockerAPI()
    config, err := api.Image("busybox", "latest").Config()
    if err != nil {
        return err
    }
    fmt.Println(config.architecture)
}
```
