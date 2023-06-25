package processor

import (
	"controller/pkg/logger"
	"controller/strategy/database"
	"controller/strategy/dict"
	e "controller/strategy/event"
	"controller/strategy/param"
	"errors"
	"fmt"
	"time"
)

type Strategy struct {
	strategy       *database.Strategy
	task           *database.Task
	fidReplication *database.FidReplication
	domains        []dict.DomainInfo
	fidRepPipeLine chan *e.Event
}

func (p *Strategy) Init(db *database.DataBase, size, num int32) {
	p.strategy = new(database.Strategy)
	p.strategy.Init(db)

	p.task = new(database.Task)
	p.task.Init(db)

	domain := new(database.Domain)
	domain.Init(db)

	var err error
	if p.domains, err = domain.Load(); err != nil {
		panic(fmt.Sprintf("domian load fail, err: %v", err.Error()))
	}

	p.fidRepPipeLine = make(chan *e.Event, size)
	for i := int32(0); i < num; i++ {
		go p.HandleFidRepEvent()
	}

}

func (p *Strategy) CreateStrategy(request *param.CreateStrategyRequest) (interface{}, error) {
	count, er := p.strategy.GetOrderStrategyCount(request.OrderId)
	if er != nil {
		return nil, er
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

	strategyInfo := p.createStrategy(request.OrderId, tasks)
	if err := p.strategy.Add(strategyInfo); err != nil {
		return nil, err
	}

	if err := p.processsFidReplication(strategyInfo.Tasks, e.FIDREP_UPSERT); err != nil {
		return nil, err
	}

	return &param.CreateStrategyResponse{
		Status: param.SUCCESS,
	}, nil
}

//注意， 后面优化可以建立策略索引，在内存中检索
func (p *Strategy) GetStrategy(request *param.GetStrategyRequset) (interface{}, error) {
	strategyInfo, err := p.strategy.GetStrategy(request.OrderId)
	if err != nil {
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

func (p *Strategy) GetDeleteStrategy(request *param.GetDeleteStrategyRequset) (interface{}, error) {
	strategyInfo, err := p.strategy.GetStrategy(request.OrderId)
	if err != nil {
		return nil, err
	}

	if err := p.processsFidReplication(strategyInfo.Tasks, e.FIDREP_DELETE); err != nil {
		return nil, err
	}

	return &param.GetDeleteStrategyResponse{
		Status:  param.SUCCESS,
		OrderId: request.OrderId,
	}, nil
}

func (p *Strategy) createStrategy(orderId string, tasks []dict.TaskInfo) *dict.StrategyInfo {
	strategyInfo := &dict.StrategyInfo{
		OrderId:    orderId,
		Tasks:      make([]*dict.Task, 0, len(tasks)),
		Desc:       "",
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	}

	regionNum := len(p.domains)
	for _, task := range tasks {

		avgMin := task.RepMin / regionNum
		avgMax := task.RepMax / regionNum

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

func (p *Strategy) processsFidReplication(tasks []*dict.Task, eventType int32) error {
	for _, task := range tasks {
		event := e.Event{
			Type: eventType,
			Data: task,
			Ret:  make(chan error),
		}
		p.AddFidRepEvent(&event)
		if err := <-event.Ret; err != nil {
			return err
		}
	}

	return nil
}

func (p *Strategy) createFidReplication(task *dict.Task) error {
	fidRepInfo := make(map[string]*dict.FidRepInfo, len(task.Reps))
	for region, rep := range task.Reps {
		fidRepInfo[region] = &dict.FidRepInfo{
			Region:       rep.Region,
			VirtualRep:   rep.VirtualRep,
			RealRep:      rep.RealRep,
			MinRep:       rep.MinRep,
			MaxRep:       rep.MaxRep,
			Status:       rep.Status,
			MinThreshold: GetFidReplicationMinThreshold(rep.Region),
			CreateTime:   time.Now().UnixMilli(),
			UpdateTime:   time.Now().UnixMilli(),
		}
	}

	err := p.fidReplication.Add(
		&dict.FidInfo{
			Fid:        task.Fid,
			Cid:        task.Cid,
			Rep:        fidRepInfo,
			Status:     dict.FIDREP_STATUS_INIT,
			CreateTime: time.Now().UnixMilli(),
			UpdateTime: time.Now().UnixMilli(),
		})

	if err != nil {
		return err
	}
	return nil
}

func (p *Strategy) updateFidReplication(fidInfo *dict.FidInfo, task *dict.Task) error {
	newRep := make(map[string]*dict.FidRepInfo, 0)
	now := time.Now().UnixMilli()
	for region, rep := range task.Reps {
		if fidRep, exists := fidInfo.Rep[region]; exists {
			newRep[region] = &dict.FidRepInfo{
				Region:       region,
				VirtualRep:   fidRep.VirtualRep,
				RealRep:      fidRep.RealRep + rep.RealRep,
				MinRep:       fidRep.MinRep + rep.MinRep,
				MaxRep:       fidRep.MaxRep + rep.MaxRep,
				Status:       fidRep.Status,
				MinThreshold: GetFidReplicationMinThreshold(region),
				CreateTime:   fidRep.CreateTime,
				UpdateTime:   now,
			}
		} else {
			newRep[region] = &dict.FidRepInfo{
				Region:       region,
				VirtualRep:   rep.VirtualRep,
				MinRep:       rep.MinRep,
				MaxRep:       rep.MaxRep,
				RealRep:      rep.RealRep,
				Status:       rep.Status,
				MinThreshold: GetFidReplicationMinThreshold(rep.Region),
				CreateTime:   now,
				UpdateTime:   now,
			}
		}
	}

	newFidInfo := &dict.FidInfo{
		Fid:        fidInfo.Fid,
		Cid:        fidInfo.Cid,
		Rep:        newRep,
		UpdateTime: now,
	}

	if err := p.fidReplication.UpdateFidInfo(newFidInfo.Fid, newFidInfo); err != nil {
		return err
	}
	return nil
}

func (p *Strategy) deleteFidReplication(fidInfo *dict.FidInfo, task *dict.Task) error {
	if fidInfo.CreateTime == fidInfo.UpdateTime {
		return p.fidReplication.RemoveFidInfo(fidInfo.Fid)
	}

	newRep := make(map[string]*dict.FidRepInfo, 0)
	now := time.Now().UnixMilli()
	for region, rep := range task.Reps {
		if fidRep, exists := fidInfo.Rep[region]; exists {
			newRep[region] = AdjustDeleteFidReplication(&dict.FidRepInfo{
				Region:       region,
				VirtualRep:   fidRep.VirtualRep,
				RealRep:      fidRep.RealRep - rep.RealRep,
				MinRep:       fidRep.MinRep - rep.MinRep,
				MaxRep:       fidRep.MaxRep - rep.MaxRep,
				Status:       fidRep.Status,
				MinThreshold: GetFidReplicationMinThreshold(region),
				CreateTime:   fidRep.CreateTime,
				UpdateTime:   now,
			})
		} else {
			logger.Errorf("delete fidrep not found, fid: %s, region: %s", fidInfo.Fid, region)
		}
	}

	newFidInfo := &dict.FidInfo{
		Fid:        fidInfo.Fid,
		Cid:        fidInfo.Cid,
		Rep:        newRep,
		UpdateTime: now,
	}

	return p.fidReplication.UpdateFidInfo(newFidInfo.Fid, newFidInfo)
}

func (p *Strategy) AddFidRepEvent(event *e.Event) {
	p.fidRepPipeLine <- event
}

func (p *Strategy) HandleFidRepEvent() {
	for {
		p.doHandleFidRepEvent(<-p.fidRepPipeLine)
	}
}

func (p *Strategy) doHandleFidRepEvent(msg *e.Event) {
	task := msg.Data.(*dict.Task)

	fidInfo, err := p.fidReplication.GetFidInfo(task.Fid)
	if err != nil {
		msg.Ret <- err
	}

	switch msg.Type {
	case e.FIDREP_UPSERT:
		if fidInfo == nil {
			msg.Ret <- p.createFidReplication(task)
		} else {
			msg.Ret <- p.updateFidReplication(fidInfo, task)
		}
	case e.FIDREP_DELETE:
		msg.Ret <- p.deleteFidReplication(fidInfo, task)
	}

}
