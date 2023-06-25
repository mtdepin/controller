package database

import (
	"controller/business/dict"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Domain struct {
	db *DataBase
}

func (p *Domain) Init(db *DataBase) {
	p.db = db
}

func (p *Domain) Load() (ret []dict.DomainInfo, err error) {
	for i := 0; i < Count; i++ {
		if ret, err = p.getDomain(bson.M{"status": bson.M{"$in": []int{dict.SUCCESS}}}); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *Domain) getDomain(query interface{}) ([]dict.DomainInfo, error) {
	records := p.db.Domain.Find(query)
	if records == nil {
		return nil, errors.New(fmt.Sprintf("getDomain ,find cond = %v fail", query))
	}
	num, err := records.Count()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("getDomain ,find cond = %v count fail", query))
	}

	domainInfo := make([]dict.DomainInfo, 0, num)
	if err = records.All(&domainInfo); err != nil {
		return nil, err
	}

	return domainInfo, nil
}

func (p *Domain) Add(info *dict.DomainInfo) (err error) {
	for i := 0; i < Count; i++ {
		if err = p.db.Domain.Insert(info); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}
