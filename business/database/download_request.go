package database

import (
	"controller/business/dict"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type DownloadRequest struct {
	db *DataBase
}

func (p *DownloadRequest) Init(db *DataBase) {
	p.db = db
}

func (p *DownloadRequest) Add(info *dict.DownloadRequestInfo) (err error) {
	for i := 0; i < Count; i++ {
		if _, err = p.db.DownloadRequest.Upsert(bson.M{"request_id": info.RequestId}, info); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}

func (p *DownloadRequest) GetDownloadRequestCount(requestId string) (count int, err error) {
	for i := 0; i < Count; i++ {
		if count, err = p.db.DownloadRequest.Find(bson.M{"request_id": requestId}).Count(); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}
