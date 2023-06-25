package database

import (
	"controller/pkg/logger"
	"controller/task_tracker/dict"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
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

func (p *FidReplication) update(fid string, value interface{}) error {
	_, err := p.db.FidReplication.Upsert(bson.M{"fid": fid}, value)
	logger.Infof("update fid: %v, value: %v", fid, value)
	return err
}

func (p *FidReplication) Load(limit int) (ret *[]dict.FidInfo, err error) {
	for i := 0; i < Count; i++ {
		if ret, err = p.getFidReplication(bson.M{"used": 0}, limit); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *FidReplication) getFidReplication(query interface{}, limit int) (*[]dict.FidInfo, error) {
	var records *mgo.Query
	if limit < 1 { //无限制
		records = p.db.FidReplication.Find(query)
	} else {
		records = p.db.FidReplication.Find(query).Limit(limit)
	}

	if records == nil {
		return nil, errors.New(fmt.Sprintf("getFidReplication ,find cond = %v fail", query))
	}
	num, err := records.Count()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("getFidReplication ,find cond = %v count fail", query))
	}

	fidInfos := make([]dict.FidInfo, 0, num)
	if err = records.All(&fidInfos); err != nil {
		return nil, err
	}

	return &fidInfos, nil
}

func (p *FidReplication) SearchFids(fids []string) (rets *[]dict.FidInfo, err error) {
	for i := 0; i < Count; i++ {
		if rets, err = p.getFidReplication(bson.M{"fid": bson.M{"$in": fids}}, 0); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *FidReplication) QueryFidByPage(lastValue int64, fieldName string, sort int, pageSize int) (ret *[]dict.FidInfo, err error) {
	for i := 0; i < Count; i++ {
		if ret, err = p.queryFidByPage(lastValue, fieldName, sort, pageSize); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *FidReplication) queryFidByPage(lastValue int64, fieldName string, sort int, pageSize int) (*[]dict.FidInfo, error) {
	sign := "+" //默认升序
	compare := "$gt"
	if sort < 0 { //降序
		sign = "-"
		compare = "$lt"
	}

	records := p.db.FidReplication.Find(bson.M{fieldName: bson.M{compare: lastValue}}).Sort(sign + fieldName).Limit(pageSize)
	if records == nil {
		return nil, errors.New(fmt.Sprintf("getFidReplication ,find lastValue: %v, fieldName:%v, sort, pageSize: %v fail", lastValue, fieldName, sort, pageSize))
	}
	num, err := records.Count()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("getFidReplication ,find lastValue: %v, fieldName:%v, sort, pageSize: %v fail", lastValue, fieldName, sort, pageSize))
	}

	fidInfos := make([]dict.FidInfo, 0, num)
	if err = records.All(&fidInfos); err != nil {
		return nil, err
	}

	return &fidInfos, nil
}
