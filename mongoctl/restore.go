package mongoctl

import (
	"gopkg.in/pipe.v2"

	"github.com/juju/errors"
	"github.com/urfave/cli"
)

var RestoreCommand = cli.Command{
	Name:  "restore",
	Usage: "restore a mongodb backup",
	Action: func(ctx *cli.Context) error {
		if err := restore(ctx); err != nil {
			logger.WithField("func", "restore").Error(err)
			return err
		}
		return nil
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
	logger.Info("exec restore")

	inDir := ctx.String("input")
	if inDir == "" {
		return errors.New("input directory is not defined")
	}

	res, err := findMaster(ctx)
	if err != nil {
		return errors.Annotate(err, "find master")
	}
	logger.Infof("mongo master ip is: %s", res.Address)
	logger.Infof("input directory is: %s", inDir)
	logger.Info("startup mongorestore")

	p := pipe.Line(
		pipe.Exec("/usr/bin/mongorestore", "-h", res.Address, "--dir", inDir),
	)

	output, err := pipe.CombinedOutput(p)
	if err != nil {
		logger.Error(err)
	}

	LogCombinedLines("restore result", logger, output)
	return nil
}
