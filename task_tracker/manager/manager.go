package manager

import (
	"controller/pkg/cache"
	"controller/task_tracker/database"
	"controller/task_tracker/dict"
	"controller/task_tracker/index"
	"controller/task_tracker/param"
	"controller/task_tracker/processor"
	"controller/task_tracker/statemachine"
	"controller/task_tracker/utils"
)

type Manager struct {
	orderIndex              *index.OrderIndex
	orderStateIndex         *index.OrderStateIndex
	stateMachine            *statemachine.StateMachine
	orderProcessor          *processor.OrderProcessor
	callbackUploadProcessor *processor.CallbackUploadProcessor
	callbackRepProcessor    *processor.CallbackRepProcessor
	callbackChargeProcessor *processor.CallbackChargeProcessor
	checkRepProcessor       *processor.CheckRepProcessor
	replicateProcessor      *processor.ReplicateProcessor
	deleteProcessor         *processor.DeleteProcessor
	chargeProcessor         *processor.ChargeProcessor
	repOrderChan            chan string
	chargeOrderChan         chan string
	orderCache              *cache.Cache
}

func (p *Manager) Init(db *database.DataBase) {
	p.repOrderChan = make(chan string, param.CHANAL_SIZE)
	p.chargeOrderChan = make(chan string, param.CHARGE_CHANAL_SIZE)

	p.orderIndex = new(index.OrderIndex)
	p.orderIndex.Init(db)

	p.orderStateIndex = new(index.OrderStateIndex)
	p.orderStateIndex.Init(db)

	p.stateMachine = new(statemachine.StateMachine)
	p.stateMachine.Init(p.orderIndex, p.orderStateIndex)

	p.orderProcessor = new(processor.OrderProcessor)
	p.orderProcessor.Init(db, p.stateMachine)

	p.callbackUploadProcessor = new(processor.CallbackUploadProcessor)
	p.callbackUploadProcessor.Init(p.stateMachine)

	p.callbackRepProcessor = new(processor.CallbackRepProcessor)
	p.callbackRepProcessor.Init(p.stateMachine)

	p.callbackChargeProcessor = new(processor.CallbackChargeProcessor)
	p.callbackChargeProcessor.Init(p.stateMachine)

	p.checkRepProcessor = new(processor.CheckRepProcessor)
	p.checkRepProcessor.Init(p.stateMachine, p.repOrderChan, p.chargeOrderChan)

	p.replicateProcessor = new(processor.ReplicateProcessor)
	p.replicateProcessor.Init(p.stateMachine, p.repOrderChan)

	p.deleteProcessor = new(processor.DeleteProcessor)
	p.deleteProcessor.Init(p.stateMachine)

	p.chargeProcessor = new(processor.ChargeProcessor)
	p.chargeProcessor.Init(p.stateMachine, p.chargeOrderChan)

	p.orderCache = new(cache.Cache)
	p.orderCache.InitCache(param.ORDER_CACHE_SIZE)
}

func (p *Manager) CreateTask(request *param.CreateTaskRequest) (interface{}, error) {
	return p.orderProcessor.CreateOrderTask(request)
}

func (p *Manager) UploadFinish(request *param.UploadFinishRequest) (interface{}, error) {
	ret, err := p.orderProcessor.UploadFinish(request)

	status, er := p.stateMachine.GetOrderStatus(request.OrderId)
	if er != nil {
		utils.Log(utils.WARN, "Manager UploadFinish stateMachine.GetOrderStatus", er.Error(), request)
		return ret, err
	}

	if status == dict.TASK_DEL_FAIL { //如果上传失败，删除订单失败，则重新投放删除到队列中。
		p.deleteProcessor.Add(request.OrderId)
	}

	/*if err == nil && request.Status == param.SUCCESS {
		status, er := p.stateMachine.GetOrderStatus(request.OrderId)
		if er != nil {
			utils.Log(utils.WARN, "Manager UploadFinish stateMachine.GetOrderStatus", er.Error(), request)
			return ret, err
		}

		if status == dict.TASK_UPLOAD_SUC { //查询
			p.checkRepProcessor.Add(request.OrderId)
		}

		if status == dict.TASK_REP_FAIL { //备份
			p.replicateProcessor.Add(request.OrderId)
		}
	}*/

	return ret, err
}

func (p *Manager) CallbackUpload(request *param.CallbackUploadRequest) (interface{}, error) {
	//幂等， 如果订单cid 已经处理了，直接返回,通过订单fid 的状态判读，做幂等。
	/*key := request.OrderId + request.Cid
	if p.orderCache.Search(key) { //存在，则返回.
		utils.Log(utils.WARN, "Manager callbackUpload order cid hava process finish,  request repeate filt", "", request)
		return param.CallbackUploadResponse{
			Status: param.SUCCESS,
		}, nil
	} else { //添加新key
		p.orderCache.Add(request.OrderId)
	}*/

	ret, err := p.callbackUploadProcessor.Process(request)

	if request.Status == param.SUCCESS {
		//幂等，去重.cache,超过5000,过滤末尾记录:  1, 2, 3, 4, 5.
		status, er := p.stateMachine.GetOrderStatus(request.OrderId)
		if er != nil {
			utils.Log(utils.WARN, "Manager callbackUpload stateMachine.GetOrderStatus", er.Error(), request)
			return ret, err
		}

		if status == dict.TASK_UPLOAD_SUC { //查询
			p.checkRepProcessor.Add(request.OrderId)
		}

		if status == dict.TASK_REP_FAIL { //备份
			p.replicateProcessor.Add(request.OrderId)
		}
	}

	return ret, err
}

func (p *Manager) CallbackRep(request *param.CallbackRepRequest) (interface{}, error) {
	rsp, err := p.callbackRepProcessor.Process(request)

	if status, err := p.stateMachine.GetOrderStatus(request.OrderId); err != nil {
		if status == dict.TASK_CHARGE_FAIL {
			p.chargeProcessor.Add(request.OrderId)
		}
	}

	return rsp, err
}

func (p *Manager) CallbackDelete(request *param.CallbackDeleteRequest) (interface{}, error) {
	//return p.callbackDeleteProcessor.Process(request)
	return nil, nil
}

func (p *Manager) CallbackCharge(request *param.CallbackChargeRequest) (interface{}, error) {
	rsp, err := p.callbackChargeProcessor.Process(request)

	if status, err := p.stateMachine.GetOrderStatus(request.OrderId); err != nil {
		if status == dict.TASK_CHARGE_FAIL {
			p.chargeProcessor.Add(request.OrderId)
		}
	}

	return rsp, err
}

func (p *Manager) DownloadFinish(request *param.DownloadFinishRequest) (interface{}, error) {
	return p.orderProcessor.DownloadFinish(request)
}
