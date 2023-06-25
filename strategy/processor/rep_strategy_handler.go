package processor

import (
	"controller/pkg/logger"
	"controller/strategy/dict"
	e "controller/strategy/event"
	"encoding/json"
	"errors"
	"fmt"
)

func (p *Strategy) createRepStrategyHandler() (interface{}, error) {
	for {
		event := <-p.fidRepEventChan
		if err := p.createFidRepStrategy(event); err != nil {
			content, _ := json.Marshal(event.Data)
			logger.Warnf("createRepStrategyHandler, createFidRepStrategy orderId:%v fail: %v, task:%v", event.OrderId, string(content), err.Error())
			event.Ret <- dict.FAIL
		} else {
			event.Ret <- dict.SUCCESS
		}
	}
}

func (p *Strategy) createFidRepStrategy(event *e.Event) error {
	task := event.Data
	key := p.prefix + task.Fid
	if err := p.mutex.Lock(key); err != nil {
		return err
	}
	defer p.mutex.UnLock(key)

	val, exist, err := p.cache.Get(task.Fid)
	if err != nil {
		return errors.New(fmt.Sprintf("get redis fail, err: %v", err.Error()))
	}

	if !exist {
		return errors.New(fmt.Sprintf("getRepStrategy,  get fid: %v from redis not eixst", task.Fid))
	}

	fidInfo := &dict.FidInfo{}
	if err := json.Unmarshal([]byte(val.(string)), fidInfo); err != nil {
		logger.Error(fmt.Sprintf("createFidRepStrategy json.Unmarshal fidInfo fail, err: %v, fidInfo :%v ", err.Error(), val.(string)))
		return err
	}

	//优化排序策略+ 权重
	return p.createOptimizeStrategy(event.OrderId, task, fidInfo)
}
