package database

import (
	"controller/business/dict"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UploadRequest struct {
	db *DataBase
}

func (p *UploadRequest) Init(db *DataBase) {
	p.db = db
}

func (p *UploadRequest) Add(info *dict.UploadRequestInfo) (err error) {
	for i := 0; i < Count; i++ {
		if _, err = p.db.OrgRequest.Upsert(bson.M{"request_id": info.RequestId}, info); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}

func (p *UploadRequest) GetOrgRequestCount(requestId string) (count int, err error) {
	for i := 0; i < Count; i++ {
		if count, err = p.db.OrgRequest.Find(bson.M{"request_id": requestId}).Count(); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}
