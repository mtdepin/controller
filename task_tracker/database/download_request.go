package database

import (
	"controller/task_tracker/dict"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type DownloadRequest struct {
	db *DataBase
}

func (p *DownloadRequest) Init(db *DataBase) {
	p.db = db
}

func (p *DownloadRequest) GetDownloadRequst(requestId string) (ret *dict.DownloadRequestInfo, err error) {
	for i := 0; i < Count; i++ {
		if ret, err = p.getDownloadRequst(requestId); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *DownloadRequest) getDownloadRequst(requestId string) (*dict.DownloadRequestInfo, error) {
	records := p.db.DownloadRequest.Find(bson.M{"request_id": requestId})
	if records == nil {
		return nil, errors.New(fmt.Sprintf("getDownloadRequst ,find request_id = %s fail", requestId))
	}
	num, err := records.Count()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("getDownloadRequst ,find request_id = %s count fail", requestId))
	}

	downloadRequestInfos := make([]dict.DownloadRequestInfo, 0, num)
	if err = records.All(&downloadRequestInfos); err != nil {
		return nil, err
	}

	if len(downloadRequestInfos) == 0 {
		return nil, errors.New(fmt.Sprintf("getDownloadRequst ,find request_id = %s no document fail", requestId))
	}

	return &downloadRequestInfos[0], nil
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
