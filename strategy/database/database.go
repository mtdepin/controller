package database

import (
	"controller/pkg/mongo"
	"gopkg.in/mgo.v2"
)

type DataBase struct {
	Session               *mgo.Session
	RepStrategyCollection *mgo.Collection
	TaskInfo              *mgo.Collection
	Domain                *mgo.Collection
	FidReplication        *mgo.Collection
}

var Db = &DataBase{}

func InitDB(url, user, pwd, dbName string, timeout int) error {
	session, err := mongo.CreateSession(url, user, pwd, dbName, timeout)
	if err != nil {
		return err
	}
	if Db.Session != nil {
		Db.Session.Close()
	}

	Db.Session = session
	db := session.DB(dbName)

	Db.RepStrategyCollection = db.C(RepStrategyCollection)
	Db.TaskInfo = db.C(TaskInfoCollection)
	Db.Domain = db.C(DomainCollection)
	Db.FidReplication = db.C(FidReplicationCollection)

	return createIndex(Db.RepStrategyCollection, "order_id")
}

func createIndex(collection *mgo.Collection, key string) error {
	if indexs, err := collection.Indexes(); err == nil {
		for i, _ := range indexs {
			if indexs[i].Key[0] == key {
				return nil
			}
		}
	}

	index := mgo.Index{
		Key:    []string{key},
		Unique: true,
	}

	return collection.EnsureIndex(index)
}
