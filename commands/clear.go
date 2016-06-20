package commands

import "github.com/codegangsta/cli"

var ClearCommand = cli.Command{
	Name:   "clear",
	Usage:  "Clear a replica set",
	Action: clearReplicaSet,
}

func clearReplicaSet(ctx *cli.Context) {

}
