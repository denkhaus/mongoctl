package mongoctl

import (
	"github.com/juju/errors"
	"github.com/urfave/cli"
)

var StatusCommand = cli.Command{
	Name:  "status",
	Usage: "print the replica status",
	Action: func(ctx *cli.Context) error {
		if err := status(ctx); err != nil {
			logger.WithField("func", "status").Error(err)
			return err
		}
		return nil
	},
}

func status(ctx *cli.Context) error {
	logger.Info("exec status")

	session, err := sessionFromCtx(ctx)
	if err != nil {
		return errors.Annotate(err, "session from ctx")
	}
	defer session.Close()

	result := &RsStatus{}
	if err := session.DB("admin").Run("replSetGetStatus", result); err != nil {
		return errors.Annotate(err, "get replica set status")
	}

	logger.Infof("Node\t\tState\t\tLast Heartbeat")
	for _, member := range result.Members {
		if member.LastHeartbeat != nil {
			logger.Infof("%s\t\t%s\t\t%v", member.Name, member.StateStr, member.LastHeartbeat)
		} else {
			logger.Infof("%s\t\t%s", member.Name, member.StateStr)
		}
	}

	return nil
}
