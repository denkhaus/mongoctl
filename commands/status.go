package commands

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"
)

var StatusCommand = cli.Command{
	Name:  "status",
	Usage: "print the replica status",
	Action: func(ctx *cli.Context) {
		if err := printStatus(ctx); err != nil {
			logger.WithField("func", "status").Error(err)
			os.Exit(1)
		}
	},
}

func printStatus(ctx *cli.Context) error {
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
