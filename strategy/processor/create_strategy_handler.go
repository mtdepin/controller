package processor

import (
	"controller/pkg/logger"
	"controller/strategy/dict"
	e "controller/strategy/event"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func (p *Strategy) createStrategyHandler() error {
	for {
		event := <-p.fidCreateEventChan
		if err := p.createFidStrategy(event); err != nil { //一个失败，本次任务的其它事件不做处理, 	f1, f2,f3.
			content, _ := json.Marshal(event.Data)
			logger.Warn("createStrategyHandler, createFidStrategy fail, task:", event.OrderId, string(content), err.Error())
			event.Ret <- dict.FAIL
		} else {
			event.Ret <- dict.SUCCESS
		}
	}
}

func (p *Strategy) createFidStrategy(event *e.Event) error {
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

	if !exist { //add new record
		return p.createInitStrategy(event.OrderId, task)
	}

	fidInfo := &dict.FidInfo{}
	if err := json.Unmarshal([]byte(val.(string)), fidInfo); err != nil {
		logger.Error(fmt.Sprintf("createFidStrategy json.Unmarshal fidInfo fail, err: %v, fidInfo :%v ", err.Error(), val.(string)))
		return err
	}

	//优化排序策略+ 权重
	return p.createOptimizeStrategy(event.OrderId, task, fidInfo)
}

func (p *Strategy) createInitStrategy(orderId string, task *dict.Task) error {
	fidRepInfo := make(map[string]map[string]*dict.Rep, 10)
	for region, rep := range task.Reps {
		mOrders := make(map[string]*dict.Rep, 5)
		p.setFidOrders(region, orderId, mOrders, rep, rep.MinRep, rep.MaxRep, rep.MinRep, rep.MaxRep)
		fidRepInfo[region] = mOrders
	}

	fidRecord := &dict.FidInfo{
		Fid:        task.Fid,
		Reps:       fidRepInfo,
		Status:     0, //init,
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	}

	bt, err := json.Marshal(fidRecord)
	if err != nil {
		logger.Error(fmt.Sprintf("createInitStrategy, json marshal fidInfo fail: %v, fidInfo: %v", err.Error(), fidRecord))
		return err
	}
	if err := p.cache.Set(task.Fid, string(bt), 0); err != nil { //task_tracker,同步设置redis, used, 还要添加lock.
		return err
	}

	if err := p.fidReplicate.Update(fidRecord.Fid, bson.M{"$set": bson.M{"reps": fidRecord.Reps, "status": fidRecord.Status, "create_time": fidRecord.CreateTime, "update_time": fidRecord.UpdateTime}}); err != nil {
		return err
	}

	return nil
}

//1. get weight.
//2. 计算minMax, maxMax.
//3. update redis, mongo.
//4. return, delete, order, region: delete,region, order.
func (p *Strategy) createOptimizeStrategy(orderId string, task *dict.Task, fidInfo *dict.FidInfo) error {
	for region, requestTask := range task.Reps { //[region] [order]{ 1 .... n}, delete, cd , order 1, order1// 10,
		if fidOrders, ok := fidInfo.Reps[region]; ok {

			initMinRep, initMaxRep, realMinRep, realMaxRep := p.estimate.CalculateRep(orderId, fidOrders, requestTask) //统一,添加到map中。
			//是最大记录且记录数大于0，才能更新real最小，最大.
			p.setFidOrders(region, orderId, fidOrders, requestTask, initMinRep, initMaxRep, realMinRep, realMaxRep) //to do test, 备份的时候，不需要添加。
		} else { //region不存在,新建region 订单
			mOrders := make(map[string]*dict.Rep, 5)
			p.setFidOrders(region, orderId, mOrders, requestTask, requestTask.MinRep, requestTask.MaxRep, requestTask.MinRep, requestTask.MaxRep)
			if len(fidInfo.Reps) == 0 {
				fidInfo.Reps = make(map[string]map[string]*dict.Rep)
			}
			fidInfo.Reps[region] = mOrders
		}
	}

	//重置used
	//fidInfo.Used = 0
	bt, err := json.Marshal(fidInfo)
	if err != nil {
		//logger.Error(fmt.Sprintf("createOptimizeStrategy, json marshal fidInfo fail: %v, fidInfo: %v", err.Error(), fidInfo))
		return err
	}
	if err := p.cache.Set(task.Fid, string(bt), 0); err != nil { //注意验证ttl, 注意可以优化点，如果值没变化，可以不用更新。
		return err
	}

	//注意只能更新备份信息，不然cid 就被覆盖了,在设置cid 的时候，没有更新到redis.
	if err := p.fidReplicate.Update(task.Fid, bson.M{"$set": bson.M{"reps": fidInfo.Reps, "used": 0}}); err != nil {
		return err
	}

	return nil
}

func (p *Strategy) setFidOrders(region, orderId string, mOrders map[string]*dict.Rep, rep *dict.RepInfo, minRep, maxRep, realMinRep, realMaxRep int) {
	if _, ok := mOrders[orderId]; !ok { //不存在，则添加，存在了，不做处理.
		mOrders[orderId] = &dict.Rep{
			Region:     region,
			MinRep:     minRep,
			MaxRep:     maxRep,
			RealRep:    rep.RealRep,
			RealMinRep: realMinRep,
			RealMaxRep: realMaxRep, //注意: 都用minRep,maxRep
			Expire:     rep.Expire,
			Status:     rep.Status,
			CreateTime: time.Now().UnixMilli(),
			UpdateTime: time.Now().UnixMilli(),
		}
	}
}
