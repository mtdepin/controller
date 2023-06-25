package database

import (
	"controller/pkg/mongo"
	"gopkg.in/mgo.v2"
)

type DataBase struct {
	Session         *mgo.Session
	OrderInfo       *mgo.Collection
	OrderState      *mgo.Collection
	OrgRequest      *mgo.Collection
	DownloadRequest *mgo.Collection
	FidReplication  *mgo.Collection
}

var Db = &DataBase{}

func InitDB(url, dbName string, timeout int) error {
	session, err := mongo.CreateSession(url, timeout)
	if err != nil {
		return err
	}
	if Db.Session != nil {
		Db.Session.Close()
	}

	Db.Session = session
	db := session.DB(dbName)

	Db.OrderInfo = db.C(OrderInfoCollection)
	Db.OrderState = db.C(OrderStatusCollection)
	Db.OrgRequest = db.C(OrgRequestCollection)
	Db.DownloadRequest = db.C(DownloadRequestCollection)
	Db.FidReplication = db.C(FidReplicationCollection)

	return nil
}
