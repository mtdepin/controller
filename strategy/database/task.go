package database

import (
	"controller/strategy/dict"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Task struct {
	db *DataBase
}

func (p *Task) Init(db *DataBase) {
	p.db = db
}

func (p *Task) GetTask(orderId string) (taskInfos []dict.TaskInfo, err error) {
	for i := 0; i < Count; i++ {
		if taskInfos, err = p.getTask(orderId); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}

func (p *Task) getTask(orderId string) ([]dict.TaskInfo, error) {
	records := p.db.TaskInfo.Find(bson.M{"order_id": orderId})
	if records == nil {
		return nil, errors.New(fmt.Sprintf("getTask ,find order_id = %s fail", orderId))
	}

	num, err := records.Count()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("getTask ,find order_id = %s count fail", orderId))
	}

	taskInfos := make([]dict.TaskInfo, 0, num)
	if err = records.All(&taskInfos); err != nil {
		return nil, err
	}

	if len(taskInfos) == 0 {
		return nil, errors.New(fmt.Sprintf("getTask ,find order_id = %s no document fail", orderId))
	}

	return taskInfos, nil
}
