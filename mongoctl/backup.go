package mongoctl

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/urfave/cli"
	"gopkg.in/pipe.v2"
)

var BackupCommand = cli.Command{
	Name:  "backup",
	Usage: "backup a mongodb",
	Action: func(ctx *cli.Context) error {
		if err := backup(ctx); err != nil {
			logger.WithField("func", "backup").Error(err)
			return err
		}
		return nil
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "output, o",
			Value:  "",
			Usage:  "output folder for backup",
			EnvVar: "MONGOCTL_BACKUP_FOLDER",
		},
		cli.StringFlag{
			Name:   "db",
			Value:  "",
			Usage:  "database to backup",
			EnvVar: "MONGOCTL_DB_TO_BACKUP",
		},
	},
}

func backup(ctx *cli.Context) error {
	logger.Info("exec backup")

	outDir := ctx.String("output")
	if outDir == "" {
		return errors.New("output directory is not defined")
	}
	
	database := ctx.String("db")
	if database == "" {
		return errors.New("database name to backup is not defined")
	}

	res, err := findMaster(ctx)
	if err != nil {
		return errors.Annotate(err, "find master")
	}

	logger.Infof("mongo master ip is: %s", res.Address)
	logger.Infof("output directory is: %s", outDir)
	logger.Info("startup mongodump")

	p := pipe.Script(
		pipe.Exec("rm", "-rf", fmt.Sprintf("%s/*", outDir)),		
		pipe.Exec("/usr/bin/mongodump", "-v", "-h", res.Address, "-o", outDir, "-d", database),
	)

	output, err := pipe.CombinedOutput(p)
	if err != nil {
		logger.Error(err)
	}

	LogCombinedLines("backup result", logger, output)
	return nil
}
