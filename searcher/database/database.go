package database

import (
	"controller/pkg/mongo"
	"gopkg.in/mgo.v2"
)

type DataBase struct {
	Session    *mgo.Session
	OrgRequest *mgo.Collection
	TaskInfo   *mgo.Collection
	Domain     *mgo.Collection
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
	return nil
}
