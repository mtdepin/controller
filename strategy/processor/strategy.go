package processor

import (
	"controller/pkg/distributionlock"
	"controller/pkg/logger"
	"controller/pkg/newcache"
	"controller/strategy/algorithm"
	"controller/strategy/config"
	"controller/strategy/database"
	"controller/strategy/dict"
	e "controller/strategy/event"
	"controller/strategy/param"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Strategy struct {
	strategy           *database.Strategy
	task               *database.Task
	domains            []*dict.DomainInfo
	fidCreateEventChan chan *e.Event
	fidDeleteEventChan chan *e.Event
	fidRepEventChan    chan *e.Event
	fidReplicate       *database.FidReplication
	estimate           *algorithm.Estimate
	cache              newcache.Cache
	mutex              *distributionlock.MutexLock
	prefix             string
	nThread            int
}

func (p *Strategy) Init(db *database.DataBase, nSize int) {
	p.strategy = new(database.Strategy)
	p.strategy.Init(db)

	p.task = new(database.Task)
	p.task.Init(db)

	p.fidReplicate = new(database.FidReplication)
	p.fidReplicate.Init(db)

	domain := new(database.Domain)
	domain.Init(db)

	p.estimate = new(algorithm.Estimate)

	domains, err := domain.Load()
	if err != nil {
		panic(fmt.Sprintf("domian load fail, err: %v", err.Error()))
	}
	nLen := len(domains)
	p.domains = make([]*dict.DomainInfo, nLen, nLen)
	for i := 0; i < nLen; i++ {
		p.domains[i] = &domains[i]
	}

	p.initRedis()

	p.nThread = nSize
	p.fidCreateEventChan = make(chan *e.Event, p.nThread)
	p.fidDeleteEventChan = make(chan *e.Event, p.nThread)
	p.fidRepEventChan = make(chan *e.Event, p.nThread)

	p.prefix = config.ServerCfg.Lock.Prefix

	for i := 0; i < p.nThread; i++ {
		go p.createStrategyHandler()
		go p.createRepStrategyHandler()
		go p.createDeleteStrategyHandler()
	}
}

func (p *Strategy) initRedis() {
	fidAddrs := strings.Split(config.ServerCfg.Redis.FidAddress, ",") //注意两个redis cluster
	lockAddrs := strings.Split(config.ServerCfg.Redis.LockAddress, ",")
	fidPassword := config.ServerCfg.Redis.FidPassword
	lockPassword := config.ServerCfg.Redis.LockPassword

	p.cache = newcache.NewRedisClient(fidAddrs, fidPassword)
	p.mutex = distributionlock.NewMutexLock(lockAddrs, lockPassword)
	logger.Info("init redis cluster success")
}

func (p *Strategy) CreateStrategy(request *param.CreateStrategyRequest) (interface{}, error) {
	count, err := p.strategy.GetOrderStrategyCount(request.OrderId)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return &param.CreateStrategyResponse{
			Status: param.SUCCESS,
		}, nil
	}

	//判断orderid 的策略是否已经存在.
	tasks, err := p.task.GetTask(request.OrderId)
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, errors.New(fmt.Sprintf("CreateStrategy, orderid: %v, task is empty ", request.OrderId))
	}
	//to do update fidreplicate , used = 0
	strategyInfo := p.createStrategy(request.OrderId, request.Region, tasks)

	if err := p.getOrderStrategy(request.OrderId, strategyInfo, p.fidCreateEventChan); err != nil {
		return nil, err
	}

	if err := p.strategy.Add(strategyInfo); err != nil {
		return nil, err
	}

	return &param.CreateStrategyResponse{
		Status: param.SUCCESS,
	}, nil
}

func (p *Strategy) GetReplicateStrategy(request *param.GetStrategyRequset) (interface{}, error) {

	t1 := time.Now().UnixMilli()
	strategyInfo, err := p.strategy.GetStrategy(request.OrderId)
	if err != nil {
		return nil, err
	}

	t2 := time.Now().UnixMilli()
	logger.Infof("GetReplicateStrategy, getOrderStrategy orderId: %v, begin: costtime: %v ms ", request.OrderId, t2-t1)
	if err := p.getOrderStrategy(request.OrderId, strategyInfo, p.fidRepEventChan); err != nil {
		return nil, err
	}

	t3 := time.Now().UnixMilli()
	logger.Infof("GetReplicateStrategy, getOrderStrategy orderId: %v, end: costtime: %v ms ", request.OrderId, t3-t2)
	//更新到strategy.
	if err := p.strategy.Update(request.OrderId, strategyInfo); err != nil {
		return nil, err
	}

	return &param.GetStrategyResponse{ //update, fidInfo
		Status:  param.SUCCESS,
		OrderId: request.OrderId,
		Strategy: &param.StrategyInfo{
			Tasks: strategyInfo.Tasks,
		},
	}, nil
}

func (p *Strategy) GetOrderDeleteStrategy(request *param.GetStrategyRequset) (interface{}, error) {
	strategyInfo, err := p.strategy.GetStrategy(request.OrderId)
	if err != nil {
		return nil, err
	}

	if err := p.getOrderStrategy(request.OrderId, strategyInfo, p.fidDeleteEventChan); err != nil {
		return nil, err
	}

	//更新到strategy.
	if err := p.strategy.Update(request.OrderId, strategyInfo); err != nil {
		return nil, err
	}

	return &param.GetStrategyResponse{
		Status:  param.SUCCESS,
		OrderId: request.OrderId,
		Strategy: &param.StrategyInfo{
			Tasks: strategyInfo.Tasks,
		},
	}, nil
}

func (p *Strategy) GetFidDeleteStrategy(request *param.GetFidDeleteStrategyRequest) (interface{}, error) {
	strategyInfo, err := p.strategy.GetStrategy(request.OrderId)
	if err != nil {
		return nil, err
	}

	rets := make([]chan int, 0, len(strategyInfo.Tasks))

	deleteTasks := make([]*dict.Task, 0, len(strategyInfo.Tasks))

	count := 0

	for _, task := range strategyInfo.Tasks {
		if _, ok := request.Fids[task.Fid]; !ok { //如果不存在则不处理。
			continue
		}

		count++

		deleteTasks = append(deleteTasks, task)
		ch := make(chan int)

		p.fidDeleteEventChan <- &e.Event{
			OrderId: request.OrderId,
			Data:    task,
			Ret:     ch,
		}
		rets = append(rets, ch)

		if count >= p.nThread-1 { //等待接收完成。
			orderStatus := dict.SUCCESS
			for _, ret := range rets {
				status := <-ret
				if status != dict.SUCCESS {
					orderStatus = dict.FAIL
				}
			}

			count = 0
			rets = rets[0:0] //清空数组
			if orderStatus == dict.FAIL {
				return nil, errors.New("get delete strategy fail")
			}
		}
	}

	//接收完rets
	orderStatus := dict.SUCCESS
	for _, ret := range rets {
		status := <-ret
		if status != dict.SUCCESS {
			orderStatus = dict.FAIL
		}
	}
	if orderStatus == dict.FAIL {
		return nil, errors.New("get delete strategy fail")
	}

	//更新到strategy.
	if err := p.strategy.Update(request.OrderId, strategyInfo); err != nil {
		return nil, err
	}

	return &param.GetStrategyResponse{
		Status:  param.SUCCESS,
		OrderId: request.OrderId,
		Strategy: &param.StrategyInfo{
			Tasks: deleteTasks,
		},
	}, nil
}

//本地上传节点备份排第一， 中心节点备份排第二， 确保数据能下载下来, 解决打洞容易失败。
func (p *Strategy) sortDomain(uploadRegion string) {
	nLen := len(p.domains)
	count := 0
	for i := 0; i < nLen; i++ {
		if p.domains[i].Level == dict.CENTER_REGION && p.domains[i].Region != uploadRegion && nLen > 1 {
			temp := p.domains[1]
			p.domains[1] = p.domains[i]
			p.domains[i] = temp
			count++
		}

		if p.domains[i].Region == uploadRegion {
			temp := p.domains[0]
			p.domains[0] = p.domains[i]
			p.domains[i] = temp
			count++
		}

		if count == 2 {
			break
		}
	}
}

func (p *Strategy) createStrategy(orderId, region string, tasks []dict.TaskInfo) *dict.StrategyInfo {
	strategyInfo := &dict.StrategyInfo{
		OrderId:    orderId,
		Tasks:      make([]*dict.Task, 0, len(tasks)),
		Desc:       "",
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	}

	p.sortDomain(region) //将本地集群排在第一个数据。

	regionNum := len(p.domains)
	if regionNum > 0 { //限制在一个集群备份。
		regionNum = 1
	}

	for _, task := range tasks {
		if task.RepMin == 0 && task.RepMax == 0 { //支持s3cmd,没有备份数，设置默认备份数。
			p.setDefaultReps(&task)
		}

		avgMin := task.RepMin / regionNum
		avgMax := task.RepMax / regionNum

		//

		strategyTask := &dict.Task{
			Fid:  task.Fid,
			Cid:  task.Cid,
			Reps: make(map[string]*dict.RepInfo, regionNum),
		}

		if task.RepMin > task.RepMax {
			task.RepMin = task.RepMax
		}

		if avgMax > 0 && avgMin > 0 {
			for i := 0; i < regionNum-1; i++ {
				rep := &dict.RepInfo{
					Region:     p.domains[i].Region,
					MinRep:     avgMin,
					MaxRep:     avgMax,
					Expire:     task.Expire,
					Encryption: task.Encryption,
				}
				strategyTask.Reps[rep.Region] = rep
			}
			//last region
			rep := &dict.RepInfo{
				Region:     p.domains[regionNum-1].Region,
				MinRep:     task.RepMin - avgMin*(regionNum-1),
				MaxRep:     task.RepMax - avgMax*(regionNum-1),
				Expire:     task.Expire,
				Encryption: task.Encryption,
			}
			strategyTask.Reps[rep.Region] = rep
			strategyInfo.Tasks = append(strategyInfo.Tasks, strategyTask)
			continue
		}

		if task.RepMin == 0 {
			rep := &dict.RepInfo{
				Region:     p.domains[0].Region,
				MinRep:     task.RepMin,
				MaxRep:     task.RepMax,
				Expire:     task.Expire,
				Encryption: task.Encryption,
			}
			strategyTask.Reps[rep.Region] = rep
			strategyInfo.Tasks = append(strategyInfo.Tasks, strategyTask)
			continue
		}

		if avgMax > 0 && avgMin == 0 && task.RepMin > 0 {
			for i := 0; i < task.RepMin-1; i++ {
				rep := &dict.RepInfo{
					Region:     p.domains[i].Region,
					MinRep:     1,
					MaxRep:     avgMax,
					Expire:     task.Expire,
					Encryption: task.Encryption,
				}
				strategyTask.Reps[rep.Region] = rep
			}
			//last region
			rep := &dict.RepInfo{
				Region:     p.domains[task.RepMin-1].Region,
				MinRep:     1,
				MaxRep:     task.RepMax - avgMax*(task.RepMin-1),
				Expire:     task.Expire,
				Encryption: task.Encryption,
			}
			strategyTask.Reps[rep.Region] = rep
			strategyInfo.Tasks = append(strategyInfo.Tasks, strategyTask)
			continue
		}

		if avgMax <= 0 {
			for i := 0; i < task.RepMin-1; i++ {
				rep := &dict.RepInfo{
					Region:     p.domains[i].Region,
					MinRep:     1,
					MaxRep:     avgMax,
					Expire:     task.Expire,
					Encryption: task.Encryption,
				}
				strategyTask.Reps[rep.Region] = rep
			}

			rep := &dict.RepInfo{
				Region:     p.domains[task.RepMin-1].Region,
				MinRep:     1,
				MaxRep:     task.RepMax - task.RepMin + 1,
				Expire:     task.Expire,
				Encryption: task.Encryption,
			}
			strategyTask.Reps[rep.Region] = rep
			strategyInfo.Tasks = append(strategyInfo.Tasks, strategyTask)
			continue
		}
	}

	return strategyInfo
}

//
func (p *Strategy) setDefaultReps(task *dict.TaskInfo) {
	switch task.Level {
	case dict.LOW:
		task.RepMin = dict.LOW_REP_NUM
		task.RepMax = dict.LOW_REP_NUM
	case dict.MIDDLE:
		task.RepMin = dict.MID_REP_NUM
		task.RepMax = dict.MID_REP_NUM
	case dict.HIGH:
		task.RepMin = dict.HIGH_REP_NUM
		task.RepMax = dict.HIGH_REP_NUM
	default:
		task.RepMin = dict.MID_REP_NUM
		task.RepMax = dict.MID_REP_NUM
	}
}

func (p *Strategy) getOrderStrategy(orderId string, strategyInfo *dict.StrategyInfo, fidEventChan chan *e.Event) error {
	nLen := len(strategyInfo.Tasks)
	div := nLen / p.nThread

	t1 := time.Now().UnixMilli()
	logger.Infof("getOrderStrategy  getPartTaskStrategy1 orderId: %v, begin: costtime: %v ms ", orderId, t1)

	for i := 0; i < div; i++ {
		if err := p.getPartTaskStrategy(i*p.nThread, (i+1)*p.nThread, orderId, strategyInfo, fidEventChan); err != nil {
			return err
		}
	}

	t2 := time.Now().UnixMilli()
	logger.Infof("getOrderStrategy  getPartTaskStrategy1 orderId: %v, end: costtime: %v ms ", orderId, t2-t1)
	//mod
	if err := p.getPartTaskStrategy(div*p.nThread, nLen, orderId, strategyInfo, fidEventChan); err != nil {
		return err
	}

	t3 := time.Now().UnixMilli()
	logger.Infof("getOrderStrategy  getPartTaskStrategy2 orderId: %v, end: costtime: %v ms ", orderId, t3-t2)

	return nil
}

//获取分片任务测试
func (p *Strategy) getPartTaskStrategy(beginPos, endPos int, orderId string, strategyInfo *dict.StrategyInfo, fidEventChan chan *e.Event) error {
	nLen := endPos - beginPos
	if nLen == 0 {
		return nil
	}

	rets := make([]chan int, 0, nLen)
	for i := beginPos; i < endPos; i++ {
		ch := make(chan int)

		fidEventChan <- &e.Event{
			OrderId: orderId,
			Data:    strategyInfo.Tasks[i],
			Ret:     ch,
		}
		rets = append(rets, ch)
	}

	orderStatus := dict.SUCCESS
	for _, ret := range rets {
		status := <-ret
		if status != dict.SUCCESS {
			orderStatus = dict.FAIL
		}
	}

	if orderStatus == dict.FAIL {
		return errors.New("create strategy fail")
	}

	return nil
}
