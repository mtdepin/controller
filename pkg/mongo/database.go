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

var mongoDb DataBase

func ConnectMgoDB(url, dbName string, timeout int) error {
	session, err := CreateSession(url, timeout)
	if err != nil {
		return err
	}
	if mongoDb.mgoSession != nil {
		mongoDb.mgoSession.Close()
	}
	mongoDb.mgoSession = session
	db := session.DB(dbName)
	mongoDb.mgoOrder = db.C("order")
	return nil
}

func CreateSession(url string, timeout int) (*mgo.Session, error) {
	dial := &mgo.DialInfo{
		Addrs:     strings.Split(url, ","),
		Direct:    true,
		Timeout:   time.Duration(timeout) * time.Second,
		PoolLimit: 100,
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
