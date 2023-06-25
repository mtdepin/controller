package processor

import (
	"controller/pkg/logger"
	"controller/strategy/dict"
	e "controller/strategy/event"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

func (p *Strategy) createDeleteStrategyHandler() (interface{}, error) {
	for {
		event := <-p.fidDeleteEventChan
		if err := p.createFidDeleteStrategy(event); err != nil {
			content, _ := json.Marshal(event.Data)
			logger.Warnf("createDeleteStrategyHandler, createFidDeleteStrategy orderId:%v fail, task:%v, err:%v", event.OrderId, content, err.Error())
			event.Ret <- dict.FAIL
		} else {
			event.Ret <- dict.SUCCESS
		}
	}
}

func (p *Strategy) createFidDeleteStrategy(event *e.Event) error {
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

	if !exist { // 都没关联，可能已经删除了,使用自己保存的删除策略.
		logger.Warnf(fmt.Sprintf("createFidDeleteStrategy, get fid: %v from redis not eixst", task.Fid))
		return nil
	}

	fidInfo := &dict.FidInfo{}
	if err := json.Unmarshal([]byte(val.(string)), fidInfo); err != nil {
		logger.Error(fmt.Sprintf("createFidDeleteStrategy json.Unmarshal fidInfo fail, err: %v, fidInfo :%v ", err.Error(), val.(string)))
		return err
	}

	//优化排序策略+ 权重
	return p.createDeleteStrategy(event.OrderId, task, fidInfo)
}

//1. 如果task.Reps 分区备份为空,region//上传区域, 删除，分为上传区域删除，备份区域删除.
//2. 计算上传区域。
//3. 计算备份区域。
func (p *Strategy) createDeleteStrategy(orderId string, task *dict.Task, fidInfo *dict.FidInfo) error {
	for region, requestRep := range task.Reps {
		if fidOrders, ok := fidInfo.Reps[region]; ok {
			if err := p.estimate.CalculateDeleteRep(orderId, fidOrders, requestRep); err != nil {
				return err
			}
			if len(fidOrders) == 0 { //此region 不在有订单。
				delete(fidInfo.Reps, region)
			}
		} else { //此region 不存在，直接返回可以删除。
			requestRep.MinRep = 0
			requestRep.MaxRep = 0
			requestRep.Status = 0
		}
	}

	//如果所有region 都删除了， 但是有一个订单关联，则，设置任意region 为删除成功，不能删除。

	if len(fidInfo.Reps) == 0 { //fid 所有集群都没备份了， 没有关联了则删除,且没有新来的订单占用
		if err := p.cache.Delete(fidInfo.Fid); err != nil {
			return err
		}
		if err := p.fidReplicate.Delete(fidInfo.Fid); err != nil {
			return err
		}
	} else { //update
		bt, err := json.Marshal(fidInfo)
		if err != nil {
			logger.Error(fmt.Sprintf("createDeleteStrategy, json marshal fidInfo fail: %v, fidInfo: %v", err.Error(), fidInfo))
			return err
		}

		if err := p.cache.Set(fidInfo.Fid, string(bt), 0); err != nil {
			return err
		}

		//注意只能更新备份信息，不然cid 就被覆盖了,在设置cid 的时候，没有更新到redis.
		if err := p.fidReplicate.Update(task.Fid, bson.M{"$set": bson.M{"reps": fidInfo.Reps}}); err != nil {
			return err
		}
	}

	return nil
}
