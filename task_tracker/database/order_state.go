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
		if ret, err = p.getOrderState(bson.M{"status": bson.M{"$in": []int{dict.TASK_INIT, dict.TASK_UPLOAD_SUC, dict.TASK_DOWNLOAD_SUC, dict.TASK_BEGIN_REP, dict.TASK_REP_SUC, dict.TASK_REP_FAIL, dict.TASK_DEL_FAIL, dict.TASK_CHARGE_FAIL}}}); err == nil {
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

func (p *OrderState) GetOrderStateByOrderId(orderId string) (ret *[]dict.OrderStateInfo, err error) {
	for i := 0; i < Count; i++ {
		if ret, err = p.getOrderState(bson.M{"order_id": orderId}); err == nil {
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

func (p *OrderState) LoadByStatus(status int) (ret *[]dict.OrderStateInfo, err error) {
	for i := 0; i < Count; i++ {
		if ret, err = p.getOrderState(bson.M{"status": status}); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *OrderState) Query(query interface{}, limit int) (ret *[]dict.OrderStateInfo, err error) {
	for i := 0; i < Count; i++ {
		if ret, err = p.query(query, limit); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *OrderState) query(query interface{}, limit int) (*[]dict.OrderStateInfo, error) {
	records := p.db.OrderState.Find(query).Limit(limit)
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

func (p *OrderState) QueryOrderByPage(orderType int32, lastValue int64, fieldName string, sort int, pageSize int) (ret *[]dict.OrderStateInfo, err error) {
	for i := 0; i < Count; i++ {
		if ret, err = p.queryOrderByPage(orderType, lastValue, fieldName, sort, pageSize); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *OrderState) queryOrderByPage(orderType int32, lastValue int64, fieldName string, sort int, pageSize int) (*[]dict.OrderStateInfo, error) {
	sign := "+" //默认升序
	compare := "$gt"
	if sort < 0 { //降序
		sign = "-"
		compare = "$lt"
	}

	records := p.db.OrderState.Find(bson.M{fieldName: bson.M{compare: lastValue}, "order_type": orderType}).Sort(sign + fieldName).Limit(pageSize)
	if records == nil {
		return nil, errors.New(fmt.Sprintf("getOrderState ,find lastValue: %v, fieldName:%v, sort, pageSize: %v fail", lastValue, fieldName, sort, pageSize))
	}
	num, err := records.Count()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("getOrderState ,find lastValue: %v, fieldName:%v, sort, pageSize: %v fail", lastValue, fieldName, sort, pageSize))
	}

	orderStateInfo := make([]dict.OrderStateInfo, 0, num)
	if err = records.All(&orderStateInfo); err != nil {
		return nil, err
	}

	return &orderStateInfo, nil
}
