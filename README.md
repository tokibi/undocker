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

```bash
undocker extract busybox:latest ./image
```

Extract directly from docker registry.

```bash
export REGISTRY_USER=xxx # optional
export REGISTRY_PASS=xxx # optional
undocker -r "https://registry-1.docker.io/" extract busybox:latest ./image
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
