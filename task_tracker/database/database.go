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
	TaskInfo        *mgo.Collection
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

	Db.OrderInfo = db.C(OrderInfoCollection)
	Db.OrderState = db.C(OrderStatusCollection)
	Db.OrgRequest = db.C(OrgRequestCollection)
	Db.DownloadRequest = db.C(DownloadRequestCollection)
	Db.FidReplication = db.C(FidReplicationCollection)
	Db.TaskInfo = db.C(TaskInfoCollection)

	orderId := "order_id"
	if err := createIndex(Db.OrderInfo, orderId); err != nil {
		return err
	}

	if err := createIndex(Db.OrderState, orderId); err != nil {
		return err
	}

	fid := "fid"
	if err := createIndex(Db.FidReplication, fid); err != nil {
		return err
	}

	return nil
}

func createIndex(collection *mgo.Collection, key string) error {
	if indexs, err := collection.Indexes(); err == nil {
		for i, _ := range indexs {
			if indexs[i].Key[0] == key {
				//fmt.Printf(" index key : %v, exist \n", key)
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
