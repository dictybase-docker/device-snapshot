package main

import (
	"os"

	"github.com/dictyBase-docker/device-snapshot/internal/command"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Name = "snapshot"
	app.Action = command.GenerateSnapshot
	app.Usage = "generate webpage snapshot using remote chrome browser"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-format",
			Usage: "format of the logging out, either of json or text.",
			Value: "json",
		},
		cli.StringFlag{
			Name:  "log-level",
			Usage: "log level for the application",
			Value: "error",
		},
		cli.StringFlag{
			Name:  "host,H",
			Usage: "remote host address",
		},
		cli.StringSliceFlag{
			Name:  "path,p",
			Usage: "webpage paths for which the snapshots will be taken",
		},
		cli.StringFlag{
			Name:   "remote-chrome-host,rh",
			Usage:  "remote chrome host",
			EnvVar: "REMOTE_CHROME_HOST",
		},
		cli.IntFlag{
			Name:   "remote-chrome-port,rp",
			Usage:  "remote chrome port",
			EnvVar: "REMOTE_CHROME_PORT",
			Value:  9222,
		},
		cli.StringFlag{
			Name:  "output,o",
			Usage: "output path for saving all the files",
		},
	}
	app.Run(os.Args)
}
