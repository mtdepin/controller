package database

import (
	"controller/task_tracker/dict"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UploadRequest struct {
	db *DataBase
}

func (p *UploadRequest) Init(db *DataBase) {
	p.db = db
}

func (p *UploadRequest) GetOrgRequest(requestId string) (ret *dict.UploadRequestInfo, err error) {
	for i := 0; i < Count; i++ {
		if ret, err = p.getOrgRequest(requestId); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *UploadRequest) getOrgRequest(requestId string) (*dict.UploadRequestInfo, error) {
	records := p.db.OrgRequest.Find(bson.M{"request_id": requestId})
	if records == nil {
		return nil, errors.New(fmt.Sprintf("GetOrgRequst ,find request_id = %s fail", requestId))
	}
	num, err := records.Count()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("GetOrgRequst ,find request_id = %s count fail", requestId))
	}

	uploadRequestInfos := make([]dict.UploadRequestInfo, 0, num)
	if err = records.All(&uploadRequestInfos); err != nil {
		return nil, err
	}

	if len(uploadRequestInfos) == 0 {
		return nil, errors.New(fmt.Sprintf("GetOrgRequst ,find request_id = %s no document fail", requestId))
	}

	return &uploadRequestInfos[0], nil
}
