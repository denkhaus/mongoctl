package mongoctl

import (
	"os"

	"gopkg.in/pipe.v2"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"
)

var RestoreCommand = cli.Command{
	Name:  "restore",
	Usage: "restore a mongodb backup",
	Action: func(ctx *cli.Context) {
		if err := restore(ctx); err != nil {
			logger.WithField("func", "restore").Error(err)
			os.Exit(1)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "input, i",
			Value:  "",
			Usage:  "input folder to restore from",
			EnvVar: "MONGOCTL_RESTORE_FOLDER",
		},
	},
}

func restore(ctx *cli.Context) error {
	inDir := ctx.String("input")
	if inDir == "" {
		return errors.New("input directory is not defined")
	}

	res, err := findMaster(ctx)
	if err != nil {
		return errors.Annotate(err, "find master")
	}

	p := pipe.Line(
		pipe.Exec("/usr/bin/mongorestore", "-h", res.Address, "--dir", inDir),
	)

	output, err := pipe.CombinedOutput(p)
	if err != nil {
		logger.Error(err)
	}

	logger.Info(output)
	return nil
}
