package commands

import (
	"os"

	"gopkg.in/mgo.v2/bson"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"
)

var AddCommand = cli.Command{
	Name:  "add",
	Usage: "Add a replica set member",
	Action: func(ctx *cli.Context) {
		if err := addMember(ctx); err != nil {
			logger.WithField("func", "add").Error(err)
			os.Exit(1)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "member, m",
			Usage: "Member <host:port> to remove",
			Value: "",
		},
		cli.BoolFlag{
			Name:  "arbitrator, a",
			Usage: "Member is arbitrator",
		},
	},
}

func addMember(ctx *cli.Context) error {
	memberHost := ctx.String("member")
	if memberHost == "" {
		return errors.New("no member host defined")
	}
	arbitrator := ctx.Bool("arbitrator")

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
	var max int64 = 0
	for _, member := range config.Members {

		if member.ID > max {
			max = member.ID
			if member.Host == memberHost {
				found = true
				break
			}
		}
	}

	if !found {
		cfg := &Host{
			ID:          max + 1,
			Host:        memberHost,
			ArbiterOnly: arbitrator,
		}

		config.Version++
		config.Members = append(config.Members, cfg)

		cmd := &bson.M{
			"replSetReconfig": config,
		}

		result := bson.M{}
		if err := session.DB("admin").Run(&cmd, &result); err != nil {
			return errors.Annotate(err, "reconfig repl set")
		}
	} else {
		return errors.Errorf("node %s already found in cluster", memberHost)
	}

	return nil
}
