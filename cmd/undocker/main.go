package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/tokibi/undocker"
	"github.com/urfave/cli"
)

var version = "unknown"

func main() {
	u := undocker.Undocker{
		Out: os.Stdout,
		Err: os.Stderr,
	}

	opts := undocker.Options{}

	app := cli.NewApp()
	app.Name = "undocker"
	app.Usage = "Decompose docker images."
	app.Version = version

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
		cli.StringFlag{
			Name:        "tmpdir",
			Value:       "/tmp/undocker",
			Usage:       "temporal directory to extract image",
			EnvVar:      "UNDOCKER_TMP_PATH",
			Destination: &opts.TmpPath,
		},
	}

	extractCommand := cli.Command{
		Name:      "extract",
		Aliases:   []string{"e"},
		Usage:     "Extract to rootfs.",
		ArgsUsage: "[image] [destination]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "overwrite-symlink-refs, s",
				Usage:       "Overwrite symbolic link references",
				Destination: &opts.Extract.OverwriteSymlinkRefs,
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
	}

	showCommand := cli.Command{
		Name:    "show",
		Aliases: []string{"s"},
		Usage:   "Show image informations",
		Subcommands: []cli.Command{
			{
				Name:      "config",
				Usage:     "Show image configuration",
				ArgsUsage: "[image]",
				Action: func(c *cli.Context) error {
					repo, tag, err := parseReference(c.Args().Get(0))
					if err != nil {
						return cli.ShowCommandHelp(c, "config")
					}
					return u.Config(repo, tag, opts)
				},
			},
			// {
			// 	Name:      "manifest",
			// 	Usage:     "Show image manifest",
			// 	ArgsUsage: "[image]",
			// 	Action: func(c *cli.Context) error {
			// 		repo, tag, err := parseReference(c.Args().Get(0))
			// 		if err != nil {
			// 			return cli.ShowCommandHelp(c, "config")
			// 		}
			// 		return u.Config(repo, tag, opts)
			// 	},
			// },
		},
	}

	app.Commands = append(app.Commands, extractCommand)
	app.Commands = append(app.Commands, showCommand)
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
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
