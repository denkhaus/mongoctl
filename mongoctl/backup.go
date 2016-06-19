package mongoctl

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	"gopkg.in/pipe.v2"
)

var BackupCommand = cli.Command{
	Name:  "backup",
	Usage: "backup a mongodb",
	Action: func(ctx *cli.Context) {
		if err := backup(ctx); err != nil {
			logger.WithField("func", "backup").Error(err)
			os.Exit(1)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "output, o",
			Value:  "",
			Usage:  "output folder for backup",
			EnvVar: "MONGOCTL_BACKUP_FOLDER",
		},
	},
}

func backup(ctx *cli.Context) error {
	logger.Info("exec backup")

	outDir := ctx.String("output")
	if outDir == "" {
		return errors.New("output directory is not defined")
	}

	res, err := findMaster(ctx)
	if err != nil {
		return errors.Annotate(err, "find master")
	}

	logger.Infof("mongo master ip is: %s", res.Address)
	logger.Infof("output directory is: %s", outDir)
	logger.Info("startup mongodump")

	p := pipe.Line(
		pipe.Exec("/usr/bin/mongodump", "-h", res.Address, "-o", outDir),
	)

	output, err := pipe.CombinedOutput(p)
	if err != nil {
		logger.Error(err)
	}

	logger.Info(string(output))
	return nil
}
