package database

import (
	"controller/task_tracker/dict"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type FidReplication struct {
	db *DataBase
}

func (p *FidReplication) Init(db *DataBase) {
	p.db = db
}

func (p *FidReplication) Search(fid string) (ret *dict.FidInfo, err error) {
	for i := 0; i < Count; i++ {
		if ret, err = p.getFidInfo(fid); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *FidReplication) getFidInfo(fid string) (*dict.FidInfo, error) {
	records := p.db.FidReplication.Find(bson.M{"fid": fid})
	if records == nil {
		return nil, errors.New(fmt.Sprintf("getDownloadRequst ,find request_id = %s fail", fid))
	}
	num, err := records.Count()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("getDownloadRequst ,find request_id = %s count fail", fid))
	}

	FidInfos := make([]dict.FidInfo, 0, num)
	if err = records.All(&FidInfos); err != nil {
		return nil, err
	}

	if len(FidInfos) == 0 {
		return nil, errors.New(fmt.Sprintf("getDownloadRequst ,find request_id = %s no document fail", fid))
	}

	return &FidInfos[0], nil
}
