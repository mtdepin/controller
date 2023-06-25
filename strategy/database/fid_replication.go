package database

import (
	"controller/strategy/dict"
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

func (p *FidReplication) Add(info *dict.FidInfo) (err error) {
	for i := 0; i < Count; i++ {
		if err = p.db.FidReplication.Insert(info); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}

func (p *FidReplication) Search(fids []string) (ret []*dict.FidInfo, err error) {
	for i := 0; i < Count; i++ {
		if ret, err = p.getFidInfos(fids); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *FidReplication) GetFidInfo(fid string) (ret *dict.FidInfo, err error) {
	for i := 0; i < Count; i++ {
		if ret, err = p.getFidInfo(fid); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *FidReplication) getFidInfos(fids []string) ([]*dict.FidInfo, error) {
	records := p.db.FidReplication.Find(bson.M{"fid": bson.M{"in": fids}})
	if records == nil {
		return nil, errors.New(fmt.Sprintf("getFidInfo ,find fid = %s fail", fids))
	}
	num, err := records.Count()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("getFidInfo ,find fid = %s count fail", fids))
	}

	FidInfos := make([]*dict.FidInfo, 0, num)
	if err = records.All(&FidInfos); err != nil {
		return nil, err
	}

	if len(FidInfos) == 0 {
		return nil, errors.New(fmt.Sprintf("getFidInfo ,find fids = %s no document fail", fids))
	}

	return FidInfos, nil
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
		return nil, errors.New(fmt.Sprintf("getDownloadRequst ,find request_id = %s no document fail", fid))
	}

	return &FidInfos[0], nil
}

func (p *FidReplication) UpdateFidInfo(fid string, info *dict.FidInfo) (err error) {
	for i := 0; i < Count; i++ {
		updateInfo := bson.M{"$set": info}
		if err = p.db.FidReplication.Update(bson.M{"fid": fid}, updateInfo); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}

func (p *FidReplication) RemoveFidInfo(fid string) (err error) {
	for i := 0; i < Count; i++ {
		if err = p.db.FidReplication.Remove(bson.M{"fid": fid}); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}
