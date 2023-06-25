package database

import (
	"controller/business/dict"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Task struct {
	db *DataBase
}

func (p *Task) Init(db *DataBase) {
	p.db = db
}

func (p *Task) Add(info *dict.TaskInfo) (err error) {
	for i := 0; i < Count; i++ {
		if err = p.db.TaskInfo.Insert(info); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}

func (p *Task) GetTaskInfo(requestId string) (ret *[]dict.TaskInfo, err error) {
	for i := 0; i < Count; i++ {
		if ret, err = p.getTaskInfo(bson.M{"request_id": requestId}); err == nil || err == mgo.ErrNotFound {
			return ret, nil
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *Task) getTaskInfo(query interface{}) (*[]dict.TaskInfo, error) {
	records := p.db.TaskInfo.Find(query)
	if records == nil {
		return nil, errors.New(fmt.Sprintf("getTaskInfo ,find cond = %v fail", query))
	}
	num, err := records.Count()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("getTaskInfo ,find cond = %v count fail", query))
	}

	taskInfos := make([]dict.TaskInfo, 0, num)
	if err = records.All(&taskInfos); err != nil {
		return nil, err
	}

	return &taskInfos, nil
}

func (p *Task) GetRequestTaskCount(requestId string) (count int, err error) {
	for i := 0; i < Count; i++ {
		if count, err = p.db.TaskInfo.Find(bson.M{"request_id": requestId}).Count(); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}

func (p *Task) GetOrderTaskCount(OrderId string) (count int, err error) {
	for i := 0; i < Count; i++ {
		if count, err = p.db.TaskInfo.Find(bson.M{"order_id": OrderId}).Count(); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}
