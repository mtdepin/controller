package database

import (
	"controller/strategy/dict"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Strategy struct {
	db *DataBase
}

func (p *Strategy) Init(db *DataBase) {
	p.db = db
}

func (p *Strategy) Add(info *dict.StrategyInfo) (err error) {
	for i := 0; i < Count; i++ {
		if err = p.db.RepStrategyCollection.Insert(info); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}

func (p *Strategy) GetStrategy(orderId string) (strategy *dict.StrategyInfo, err error) {
	for i := 0; i < Count; i++ {

		if strategy, err = p.getStrategy(orderId); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}

	return
}

func (p *Strategy) getStrategy(orderId string) (*dict.StrategyInfo, error) {
	records := p.db.RepStrategyCollection.Find(bson.M{"order_id": orderId})
	if records == nil {
		return nil, errors.New(fmt.Sprintf("getStrategy ,find orderId = %s fail", orderId))
	}

	num, err := records.Count()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("getStrategy ,find orderId = %s count fail", orderId))
	}

	strategyInfos := make([]dict.StrategyInfo, 0, num)
	if err = records.All(&strategyInfos); err != nil {
		return nil, err
	}

	if len(strategyInfos) == 0 {
		return nil, errors.New(fmt.Sprintf("getStrategy ,find orderId = %s no document fail", orderId))
	}

	return &strategyInfos[0], nil
}

func (p *Strategy) GetOrderStrategyCount(orderId string) (count int, err error) {
	for i := 0; i < Count; i++ {
		if count, err = p.db.RepStrategyCollection.Find(bson.M{"order_id": orderId}).Count(); err == nil {
			return
		}
		time.Sleep(time.Duration(TimeInternal) * time.Millisecond)
	}
	return
}
