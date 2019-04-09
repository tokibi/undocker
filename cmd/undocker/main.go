package main

import (
	"errors"
	"os"
	"strings"

	"github.com/tokibi/undocker"
	"github.com/urfave/cli"
)

func main() {
	u := undocker.Undocker{
		Out: os.Stdout,
		Err: os.Stderr,
	}

	opts := undocker.Options{}

	app := cli.NewApp()
	app.Name = "undocker"
	app.Usage = "Decompose docker images."
	app.Version = "0.1.2"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "registry-url, r",
			Usage:       "docker registry url",
			EnvVar:      "REGISTRY_URL",
			Destination: &opts.RegistryURL,
		},
		cli.StringFlag{
			Name:        "registry-user, u",
			Usage:       "docker registry login username",
			EnvVar:      "REGISTRY_USER",
			Destination: &opts.RegistryUser,
		},
		cli.StringFlag{
			Name:        "registry-pass, p",
			Usage:       "docker registry login password",
			EnvVar:      "REGISTRY_PASS",
			Destination: &opts.RegistryPass,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "extract",
			Aliases:   []string{"e"},
			Usage:     "Extract to rootfs.",
			ArgsUsage: "[image] [destination]",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "keep-symlink-refs, k",
					Usage:       "Keep the destination to which symbolic link refers",
					Destination: &opts.Extract.KeepSymlinkRefs,
				},
			},
			Action: func(c *cli.Context) error {
				repo, tag, err := parseReference(c.Args().Get(0))
				if err != nil {
					return cli.ShowCommandHelp(c, "extract")
				}

				dest := c.Args().Get(1)
				if dest == "" {
					dest = "."
				}
				return u.Extract(repo, tag, dest, opts)
			},
		},
		{
			Name:      "config",
			Aliases:   []string{"c"},
			Usage:     "Show image configuration.",
			ArgsUsage: "[image]",
			Action: func(c *cli.Context) error {
				repo, tag, err := parseReference(c.Args().Get(0))
				if err != nil {
					return cli.ShowCommandHelp(c, "config")
				}
				return u.Config(repo, tag, opts)
			},
		},
	}

	app.Run(os.Args)
}

func parseReference(arg string) (repository, tag string, err error) {
	ref := strings.SplitN(arg, ":", 2)
	if ref[0] == "" {
		return "", "", errors.New("Invalid image")
	}
	if len(ref) < 2 {
		return ref[0], "latest", nil
	}
	return ref[0], ref[1], nil
}
