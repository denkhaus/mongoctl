package main

import (
	"fmt"
	"os"

	"github.com/denkhaus/mongoctl/mongoctl"
	"github.com/sirupsen/logrus"

	"github.com/codegangsta/cli"
)

var (
	AppVersion = "0.1.0"
	Revision   = "undefined"
)
var (
	logger = logrus.WithField("pkg", "main")
)

func main() {
	app := cli.NewApp()
	app.Name = "mongoctl"
	app.EnableBashCompletion = true
	app.Version = fmt.Sprintf("%s-%s", AppVersion, Revision)
	app.Usage = "A mongo tool cli"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host, d",
			Value:  "127.0.0.1:27017",
			Usage:  "db host to work on",
			EnvVar: "MONGOCTL_DBHOST",
		},
		cli.StringFlag{
			Name:   "username, u",
			Value:  "",
			Usage:  "db username",
			EnvVar: "MONGOCTL_DBUSERNAME",
		},
		cli.StringFlag{
			Name:   "password, p",
			Value:  "",
			Usage:  "db password",
			EnvVar: "MONGOCTL_DBPASSWORD",
		},
		cli.BoolFlag{
			Name:  "revision, r",
			Usage: "Print revision",
		},
	}

	app.Action = func(ctx *cli.Context) {
		if ctx.GlobalBool("revision") {
			fmt.Println(Revision)
			return
		}

		cli.ShowAppHelp(ctx)
	}
	app.Commands = []cli.Command{
		mongoctl.StatusCommand,
		mongoctl.AddCommand,
		mongoctl.RemoveCommand,
		mongoctl.ClearCommand,
		mongoctl.BackupCommand,
		mongoctl.RestoreCommand,
	}

	logger.Infof("startup mongoctl %s", app.Version)
	app.Run(os.Args)
}
