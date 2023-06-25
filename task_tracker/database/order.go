package database

import (
	"controller/task_tracker/dict"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Order struct {
	db *DataBase
}

func (p *Order) Init(db *DataBase) {
	p.db = db
}

func (p *Order) Add(info *dict.OrderInfo) (err error) {
	for i := 0; i < Count; i++ {
		if err = p.db.OrderInfo.Insert(info); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}

func (p *Order) Load() (ret *[]dict.OrderInfo, err error) {
	for i := 0; i < Count; i++ {
		//bson.M{"status": bson.M{"$in": []int{dict.TASK_INIT, dict.TASK_UPLOAD_SUC, dict.TASK_DOWNLOAD_SUC, dict.TASK_REP_SUC, dict.TASK_REP_FAIL, dict.TASK_DEL_FAIL, dict.TASK_CHARGE_FAIL}}}
		if ret, err = p.getOrderInfo(bson.M{"status": bson.M{"$in": []int{dict.TASK_INIT, dict.TASK_UPLOAD_SUC, dict.TASK_DOWNLOAD_SUC, dict.TASK_REP_SUC, dict.TASK_REP_FAIL, dict.TASK_DEL_FAIL, dict.TASK_CHARGE_FAIL}}}); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *Order) Update(info *dict.OrderInfo) (err error) {
	for i := 0; i < Count; i++ {
		if _, err = p.db.OrderInfo.Upsert(bson.M{"order_id": info.OrderId}, info); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}

func (p *Order) GetOrder(requestId string) (ret *dict.OrderInfo, err error) {
	for i := 0; i < Count; i++ {
		if ret, err = p.getOrder(requestId); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *Order) getOrder(requestId string) (*dict.OrderInfo, error) {
	records := p.db.OrderInfo.Find(bson.M{"request_id": requestId})
	if records == nil {
		return nil, errors.New(fmt.Sprintf("getOrder ,find request_id = %s fail", requestId))
	}

	num, err := records.Count()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("getOrder ,find request_id = %s count fail", requestId))
	}

	orderInfos := make([]dict.OrderInfo, 0, num)
	if err = records.All(&orderInfos); err != nil {
		return nil, err
	}

	if len(orderInfos) == 0 {
		return nil, errors.New(fmt.Sprintf("getOrder ,find request_id = %s no document fail", requestId))
	}

	return &orderInfos[0], nil
}

func (p *Order) getOrderInfo(query interface{}) (*[]dict.OrderInfo, error) {
	records := p.db.OrderInfo.Find(query)
	if records == nil {
		return nil, errors.New(fmt.Sprintf("getOrderInfo ,find cond = %v fail", query))
	}
	num, err := records.Count()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("getOrderInfo ,find cond = %v count fail", query))
	}

	orderInfo := make([]dict.OrderInfo, 0, num)
	if err = records.All(&orderInfo); err != nil {
		return nil, err
	}

	return &orderInfo, nil
}
