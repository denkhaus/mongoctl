package mongoctl

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	"gopkg.in/mgo.v2/bson"
)

var RemoveCommand = cli.Command{
	Name:  "remove",
	Usage: "Remove a replica set member",
	Action: func(ctx *cli.Context) {
		if err := removeMember(ctx); err != nil {
			logger.WithField("func", "remove").Error(err)
			os.Exit(1)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "member, m",
			Usage: "Member <host:port> to remove",
			Value: "",
		},
	},
}

func removeMember(ctx *cli.Context) error {
	memberHost := ctx.String("member")
	if memberHost == "" {
		return errors.New("no member host defined")
	}

	session, err := sessionFromCtx(ctx)
	if err != nil {
		return errors.Annotate(err, "session from ctx")
	}
	defer session.Close()

	config := &RsConf{}
	coll := session.DB("local").C("system.replset")
	count, err := coll.Count()
	if err != nil {
		return errors.Annotate(err, "count")
	}

	if count > 1 {
		return errors.New("local.system.replset has unexpected contents")
	}

	err = coll.Find(bson.M{}).One(&config)
	if err != nil {
		return errors.Annotate(err, "get replset config")

	}

	found := false
	for i, member := range config.Members {
		if member.Host == memberHost {
			config.Members = append(config.Members[:i], config.Members[i+1:]...)
			found = true
			break
		}
	}

	if found {
		config.Version++
		cmd := &bson.M{
			"replSetReconfig": config,
		}
		result := bson.M{}
		if err := session.DB("admin").Run(&cmd, &result); err != nil {
			return errors.Annotate(err, "reconfig repl set")
		}
	} else {
		return errors.Errorf("node %s not found in cluster", memberHost)
	}

	return nil
}
