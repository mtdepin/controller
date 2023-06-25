package mongo

import (
	"gopkg.in/mgo.v2"
	"strings"
	"time"
)

type DataBase struct {
	mgoSession *mgo.Session
	mgoOrder   *mgo.Collection
}

func CreateSession(url, user, pwd, dbName string, timeout int) (*mgo.Session, error) {
	dial := &mgo.DialInfo{
		Addrs:     strings.Split(url, ","),
		Direct:    true,
		Timeout:   time.Duration(timeout) * time.Second,
		PoolLimit: 500,
		//Username:  user,
		//Password:  pwd,
		Database: dbName,
	}

	session, err := mgo.DialWithInfo(dial)
	if err != nil {
		return nil, err
	}

	session.SetSyncTimeout(time.Duration(timeout) * time.Second)
	session.SetSocketTimeout(time.Duration(timeout) * time.Second)
	session.SetMode(mgo.Eventual, true)
	return session, nil
}
