package database

import (
	"controller/task_tracker/dict"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type OrderState struct {
	db *DataBase
}

func (p *OrderState) Init(db *DataBase) {
	p.db = db
}

func (p *OrderState) Load() (ret *[]dict.OrderStateInfo, err error) {
	for i := 0; i < Count; i++ {
		if ret, err = p.getOrderState(bson.M{"status": bson.M{"$in": []int{dict.TASK_INIT, dict.TASK_UPLOAD_SUC, dict.TASK_DOWNLOAD_SUC, dict.TASK_REP_SUC, dict.TASK_REP_FAIL, dict.TASK_DEL_FAIL, dict.TASK_CHARGE_FAIL}}}); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *OrderState) Add(info *dict.OrderStateInfo) (err error) {
	for i := 0; i < Count; i++ {
		if err = p.db.OrderState.Insert(info); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}

func (p *OrderState) Update(info *dict.OrderStateInfo) (err error) {
	for i := 0; i < Count; i++ {
		if _, err = p.db.OrderState.Upsert(bson.M{"order_id": info.OrderId}, info); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}

func (p *OrderState) GetOrderState(status int) (ret *[]dict.OrderStateInfo, err error) {
	for i := 0; i < Count; i++ {
		if ret, err = p.getOrderState(bson.M{"status": status}); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *OrderState) getOrderState(query interface{}) (*[]dict.OrderStateInfo, error) {
	records := p.db.OrderState.Find(query)
	if records == nil {
		return nil, errors.New(fmt.Sprintf("getOrderState ,find cond = %v fail", query))
	}
	num, err := records.Count()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("getOrderState ,find cond = %v count fail", query))
	}

	orderStateInfo := make([]dict.OrderStateInfo, 0, num)
	if err = records.All(&orderStateInfo); err != nil {
		return nil, err
	}

	return &orderStateInfo, nil
}
