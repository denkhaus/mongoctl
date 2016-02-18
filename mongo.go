package mongoctl

import (
	"time"

	"github.com/codegangsta/cli"

	"gopkg.in/mgo.v2"
)

func sessionFromCtx(ctx *cli.Context) (*mgo.Session, error) {
	host := ctx.GlobalString("host")

	info := &mgo.DialInfo{
		Addrs:    []string{host},
		Timeout:  5 * time.Second,
		Username: ctx.GlobalString("username"),
		Password: ctx.GlobalString("password"),
	}

	return mgo.DialWithInfo(info)
}
