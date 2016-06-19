package mongoctl

import (
	"fmt"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/juju/errors"
	"github.com/urfave/cli"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	errMasterNotFound = errors.New("master node is not found")
)

func sessionFromHost(host, username, password string) (*mgo.Session, error) {
	info := &mgo.DialInfo{
		Addrs:    []string{host},
		Timeout:  5 * time.Second,
		Username: username,
		Password: password,
	}

	return mgo.DialWithInfo(info)
}

func sessionFromCtx(ctx *cli.Context) (*mgo.Session, error) {
	host := ctx.GlobalString("host")
	username := ctx.GlobalString("username")
	password := ctx.GlobalString("password")
	return sessionFromHost(host, username, password)
}

// IsMaster returns information about the configuration of the node that
// the given session is connected to.
func isMaster(session *mgo.Session) (*IsMasterResults, error) {
	results := &IsMasterResults{}

	if err := session.Run("isMaster", results); err != nil {
		return nil, err
	}

	results.Address = unFixIpv6Address(results.Address)
	results.PrimaryAddress = unFixIpv6Address(results.PrimaryAddress)
	for index, address := range results.Addresses {
		results.Addresses[index] = unFixIpv6Address(address)
	}
	return results, nil
}

// finds master in replset
func findMaster(ctx *cli.Context) (*IsMasterResults, error) {
	username := ctx.GlobalString("username")
	password := ctx.GlobalString("password")

	session, err := sessionFromCtx(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "session from ctx")
	}
	defer session.Close()

	config := &RsConf{}
	coll := session.DB("local").C("system.replset")
	count, err := coll.Count()
	if err != nil {
		return nil, errors.Annotate(err, "count")
	}

	if count > 1 {
		return nil, errors.New("local.system.replset has unexpected contents")
	}

	err = coll.Find(bson.M{}).One(&config)
	if err != nil {
		return nil, errors.Annotate(err, "get replset config")
	}

	for _, member := range config.Members {
		spew.Dump(member)
		sess, err := sessionFromHost(member.Host, username, password)
		defer sess.Close()

		if err != nil {
			return nil, errors.Annotate(err, "session from host")
		}

		res, err := isMaster(sess)
		if err != nil {
			return nil, errors.Annotate(err, "isMaster")
		}

		if res.IsMaster {
			return res, nil
		}
	}

	return nil, errMasterNotFound
}

// Turn normal ipv6 addresses into the "bad format" that mongo requires us
// to use. (Mongo can't parse square brackets in ipv6 addresses.)
func fixIpv6Address(address string) string {
	address = strings.Replace(address, "[", "", 1)
	address = strings.Replace(address, "]", "", 1)
	return address
}

// Turn "bad format" ipv6 addresses ("::1:port"), that mongo uses,  into good
// format addresses ("[::1]:port").
func unFixIpv6Address(address string) string {
	if strings.Count(address, ":") >= 2 && strings.Count(address, "[") == 0 {
		lastColon := strings.LastIndex(address, ":")
		host := address[:lastColon]
		port := address[lastColon+1:]
		return fmt.Sprintf("[%s]:%s", host, port)
	}
	return address
}
