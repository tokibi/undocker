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
	app.Usage = "Parsing docker images."
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "registry-url, r",
			Usage: "Docker registry url",
		},
		cli.StringFlag{
			Name:  "registry-user, u",
			Usage: "Docker registry username",
		},
		cli.StringFlag{
			Name:  "registry-pass, p",
			Usage: "Docker registry password",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "extract",
			Aliases: []string{"e"},
			Usage:   "Extract to rootfs.",
			Action: func(c *cli.Context) error {
				return undocker.Extract(c)
			},
		},
		// {
		// 	Name:    "config",
		// 	Aliases: []string{"c"},
		// 	Usage:   "Show image configuration.",
		// 	Action: func(c *cli.Context) error {
		// 		return nil
		// 	},
		// },
	}

	app.Run(os.Args)
}
