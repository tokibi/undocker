package main

import (
	"os"

	"github.com/tokibi/undocker"
	"github.com/urfave/cli"
)

func main() {
	undocker := undocker.Undocker{
		Out: os.Stdout,
		Err: os.Stderr,
	}

	app := cli.NewApp()
	app.Name = "undocker"
	app.Usage = "Decompose docker images."
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "registry-url, r",
			Usage:  "docker registry url",
			EnvVar: "REGISTRY_URL",
		},
		cli.StringFlag{
			Name:   "registry-user, u",
			Usage:  "docker registry login username",
			EnvVar: "REGISTRY_USER",
		},
		cli.StringFlag{
			Name:   "registry-pass, p",
			Usage:  "docker registry login password",
			EnvVar: "REGISTRY_PASS",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "extract",
			Aliases:   []string{"e"},
			Usage:     "Extract to rootfs.",
			ArgsUsage: "[image] [destination]",
			Action: func(c *cli.Context) error {
				if c.NArg() < 2 {
					return cli.ShowCommandHelp(c, "extract")
				}
				return undocker.Extract(c)
			},
		},
		{
			Name:      "config",
			Aliases:   []string{"c"},
			Usage:     "Show image configuration.",
			ArgsUsage: "[image]",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					return cli.ShowCommandHelp(c, "config")
				}
				return undocker.Config(c)
			},
		},
	}

	app.Run(os.Args)
}
