package mongoctl

import "github.com/urfave/cli"

var ClearCommand = cli.Command{
	Name:  "clear",
	Usage: "Clear a replica set",
	Action: func(ctx *cli.Context) error {
		if err := clearReplicaSet(ctx); err != nil {
			logger.WithField("func", "clear").Error(err)
			return err
		}
		return nil
	},
}

func clearReplicaSet(ctx *cli.Context) error {
	logger.Info("exec clear")

	return nil
}
