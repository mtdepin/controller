package database

import (
	"controller/pkg/mongo"
	"gopkg.in/mgo.v2"
)

type DataBase struct {
	Session         *mgo.Session
	OrgRequest      *mgo.Collection
	TaskInfo        *mgo.Collection
	Domain          *mgo.Collection
	DownloadRequest *mgo.Collection
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

	Db.TaskInfo = db.C(TaskInfoCollection)
	Db.OrgRequest = db.C(OrgRequestCollection)
	Db.Domain = db.C(DomainCollection)
	Db.DownloadRequest = db.C(DownloadRequestCollection)
	return nil
}
