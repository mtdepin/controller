package database

import (
	"controller/business/dict"
	"controller/pkg/logger"
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

func (p *FidReplication) Update(fid string, value interface{}) (err error) {
	for i := 0; i < Count; i++ {
		if err = p.update(fid, value); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *FidReplication) Delete(fid string) (err error) {
	for i := 0; i < Count; i++ {
		if err = p.db.FidReplication.Remove(bson.M{"fid": fid}); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *FidReplication) getFidInfo(fid string) (*dict.FidInfo, error) {
	records := p.db.FidReplication.Find(bson.M{"fid": fid})
	if records == nil {
		return nil, errors.New(fmt.Sprintf("getFidInfo ,find fid = %s fail", fid))
	}
	num, err := records.Count()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("getFidInfo ,find fid = %s count fail", fid))
	}

	FidInfos := make([]dict.FidInfo, 0, num)

	if err = records.All(&FidInfos); err != nil {
		return nil, err
	}

	if len(FidInfos) == 0 {
		return nil, errors.New(fmt.Sprintf("getFidInfo ,find fid = %s no document fail", fid))
	}

	return &FidInfos[0], nil
}

func (p *FidReplication) update(fid string, value interface{}) error {
	_, err := p.db.FidReplication.Upsert(bson.M{"fid": fid}, value)
	logger.Infof("update fid: %v, value: %v", fid, value)
	return err
}
