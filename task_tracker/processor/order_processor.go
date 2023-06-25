package processor

import (
	"controller/api"
	"controller/pkg/logger"
	"controller/task_tracker/database"
	"controller/task_tracker/dict"
	"controller/task_tracker/event"
	e "controller/task_tracker/event"
	"controller/task_tracker/param"
	"controller/task_tracker/statemachine"
	"controller/task_tracker/utils"
	"errors"
	"fmt"
	"time"
)

type OrderProcessor struct {
	order           *database.Order
	uploadRequest   *database.UploadRequest
	downloadRequest *database.DownloadRequest
	fidReplicate    *database.FidReplication
	stateMachine    *statemachine.StateMachine
}

func (p *OrderProcessor) Init(db *database.DataBase, machine *statemachine.StateMachine, fidReplicate *database.FidReplication) {
	p.order = new(database.Order)
	p.order.Init(db)

	p.uploadRequest = new(database.UploadRequest)
	p.uploadRequest.Init(db)

	p.downloadRequest = new(database.DownloadRequest)
	p.downloadRequest.Init(db)

	p.fidReplicate = fidReplicate
	p.stateMachine = machine
}

func (p *OrderProcessor) CreateOrderTask(request *api.CreateTaskRequest) (interface{}, error) {
	if request.Type == param.UPLOAD {
		return p.createOrderUploadTask(request)
	} else if request.Type == param.DOWNLOAD {
		return p.createOrderDownloadTask(request)
	}

	return nil, errors.New(fmt.Sprintf("unknown order type: %d ", request.Type))
}

func (p *OrderProcessor) createOrderUploadTask(request *api.CreateTaskRequest) (interface{}, error) {
	uploadRequest, err := p.uploadRequest.GetOrgRequest(request.RequestId)
	if err != nil {
		return nil, err
	}

	if uploadRequest.PieceNum == 0 {
		return nil, errors.New("piecenum == 0")
	}

	//判断订单是否创建完成
	if orderInfo, err := p.stateMachine.GetOrderByRequestId(request.RequestId); err == nil {
		return &api.CreateTaskResponse{
			Status:  param.SUCCESS,
			OrderId: orderInfo.OrderId,
		}, nil
	}

	orderEvent := p.createUploadOrderEvent(uploadRequest)

	if err := p.stateMachine.Send(orderEvent.OrderId, orderEvent); err != nil {
		return nil, err
	}

	status := <-orderEvent.Ret

	return &api.CreateTaskResponse{
		Status:  status,
		OrderId: orderEvent.OrderId,
	}, nil
}

func (p *OrderProcessor) createOrderDownloadTask(request *api.CreateTaskRequest) (interface{}, error) {
	downloadRequest, err := p.downloadRequest.GetDownloadRequst(request.RequestId)
	if err != nil {
		return nil, err
	}

	if len(downloadRequest.Tasks) == 0 {
		return nil, errors.New("download task is empty")
	}

	//判断订单是否创建完成
	if orderInfo, err := p.stateMachine.GetOrderByRequestId(request.RequestId); err == nil {
		return param.DownloadTaskResponse{
			Status:  param.SUCCESS,
			OrderId: orderInfo.OrderId,
		}, nil
	}

	orderInfo := p.generateOrder(request)

	orderEvent, err1 := p.generateDownloadOrderEvent(orderInfo.OrderId, downloadRequest)
	if err1 != nil {
		return nil, err1
	}

	if err := p.stateMachine.Send(orderInfo.OrderId, orderEvent); err != nil {
		return nil, err
	}

	status := <-orderEvent.Ret

	return param.DownloadTaskResponse{
		Status:  status,
		OrderId: orderInfo.OrderId,
	}, nil
}

func (p *OrderProcessor) UploadFinish(request *param.UploadFinishRequest) (interface{}, error) {
	status := param.FAIL
	if request.Status == param.SUCCESS {
		//time.Sleep(5 * time.Second)
		state, err := p.stateMachine.GetOrderStateInfo(request.OrderId)
		if err != nil { //订单不存在，已经处理完成。
			return &param.UploadFinishResponse{
				OrderId: request.OrderId,
				Status:  param.SUCCESS, //上传失败，删除成功，返回成功，否则返回失败。
			}, nil
		}

		if state.Status >= dict.TASK_UPLOAD_SUC {
			status = param.SUCCESS
		} else { //初始化状态，如果都是重复文件， 已经上传完成。
			if err := p.InformReplication(state); err != nil {
				utils.Log(utils.WARN, "OrderProcessor UploadFinish repeat order  InformReplication fail", err.Error(), state)
				return &param.UploadFinishResponse{
					OrderId: request.OrderId,
					Status:  param.FAIL,
				}, nil
			}

			status = param.PROCEED
		}

		return &param.UploadFinishResponse{
			OrderId: request.OrderId,
			Status:  status,
		}, nil
	}
	//处理失败,删除订单.
	event := p.generateOrderUploadFinishEvent(request.OrderId, status)

	if err := p.stateMachine.Send(request.OrderId, event); err != nil {
		return nil, err
	}

	retStatus := <-event.Ret

	//如果是上传失败，且任务删除成功，将订单从状态机中删除。
	if status == param.FAIL && retStatus == param.SUCCESS {
		p.stateMachine.Delete(request.OrderId)
	}

	return &param.UploadFinishResponse{
		OrderId: request.OrderId,
		Status:  retStatus, //上传失败，删除成功，返回成功，否则返回失败。
	}, nil
}

func (p *OrderProcessor) SetRepFailState(state *dict.OrderStateInfo) error {
	for _, task := range state.Tasks {
		if task.Repeate != dict.REPEATE { //不是重复文件，则直接返回。
			return nil
		}
	}

	for _, task := range state.Tasks {
		for _, rep := range task.Reps {
			rep.Status = dict.TASK_REP_FAIL
		}
		task.Status = dict.TASK_REP_FAIL
	}
	state.Status = dict.TASK_REP_FAIL

	return p.stateMachine.UpdateOrderStateInfo(state.OrderId, state)
}

func (p *OrderProcessor) generateCallbackUploadEvent(state *dict.OrderStateInfo) (*e.Event, error) {
	for _, task := range state.Tasks {
		callbackUploadEvent := &e.CallbackUploadEvent{
			OrderId: state.OrderId,
			Fid:     task.Fid,
			Cid:     task.Cid,
			Region:  task.Region,
			Origins: task.Origins,
			Status:  dict.SUCESS, //上传成功,触发文件备份。
		}

		return &e.Event{
			Type:    e.CALLBACK_UPLOAD,
			OrderId: state.OrderId,
			Ret:     make(chan int),
			Data:    callbackUploadEvent,
		}, nil
	}

	return nil, errors.New("generate callbackUploadEvent fail")

}

func (p *OrderProcessor) generateOrder(request *api.CreateTaskRequest) *dict.OrderInfo {
	return &dict.OrderInfo{
		OrderId:    CreateOrderId(),
		RequestId:  request.RequestId,
		OrderType:  request.Type,
		Status:     dict.TASK_INIT,
		Desc:       "",
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli()}
}
func (p *OrderProcessor) createUploadOrderEvent(uploadInfo *dict.UploadRequestInfo) *event.Event {
	createOrderEvent := &event.CreateOrderEvent{RequestId: uploadInfo.RequestId, OrderType: param.UPLOAD, PieceNum: uploadInfo.PieceNum}
	return &event.Event{
		Type:    event.CREATE_ORDER,
		OrderId: CreateOrderId(),
		Ret:     make(chan int),
		Data:    createOrderEvent,
	}
}

func (p *OrderProcessor) generateDownloadOrderEvent(orderId string, request *dict.DownloadRequestInfo) (*event.Event, error) {
	createOrderEvent := &event.CreateOrderEvent{RequestId: request.RequestId, OrderType: param.DOWNLOAD}

	createOrderEvent.Cids = make([]string, 0, len(request.Tasks))
	for _, task := range request.Tasks {
		createOrderEvent.Cids = append(createOrderEvent.Cids, task.Cid)
	}

	return &event.Event{
		Type:    event.CREATE_ORDER,
		OrderId: orderId,
		Ret:     make(chan int),
		Data:    createOrderEvent,
	}, nil
}

func (p *OrderProcessor) generateOrderUploadFinishEvent(orderId string, status int) *event.Event {
	return &event.Event{
		Type:    event.TASK_UPLOAD_FINISH,
		OrderId: orderId,
		Ret:     make(chan int),
		Data: &event.OrderUploadFinishEvent{
			OrderId: orderId,
			Status:  status,
		},
	}
}

func (p *OrderProcessor) DownloadFinish(request *param.DownloadFinishRequest) (interface{}, error) {
	event := p.generateOrderDownloadFinishEvent(request.OrderId, request.Status)

	if err := p.stateMachine.Send(request.OrderId, event); err != nil {
		return nil, err
	}

	retStatus := <-event.Ret

	return &param.DownloadFinishResponse{
		OrderId: request.OrderId,
		Status:  retStatus,
	}, nil
}

func (p *OrderProcessor) generateOrderDownloadFinishEvent(orderId string, status int) *event.Event {
	return &event.Event{
		Type:    event.TASK_DOWNLOAD_FINISH,
		OrderId: orderId,
		Ret:     make(chan int),
		Data: &event.OrderDownloadFinishEvent{
			OrderId: orderId,
			Status:  status,
		},
	}
}

func (p *OrderProcessor) InformReplication(state *dict.OrderStateInfo) error {
	for _, task := range state.Tasks {
		if task.Repeate != dict.REPEATE {
			return nil
		}
	}

	t1 := time.Now().UnixMilli()
	logger.Infof("InformReplication order begin: %v ", state.OrderId)
	//全部为重复订单，生成上传通知事件，触发备份。
	event, err := p.generateCallbackUploadEvent(state)
	if err != nil {
		return err
	}

	if err := p.stateMachine.Send(event.OrderId, event); err != nil {
		return err
	}

	status := <-event.Ret

	t2 := time.Now().UnixMilli()
	logger.Infof("InformReplication order end: %v, total_time: %v ", state.OrderId, t2-t1)
	if status != statemachine.SUCCESS {
		return errors.New("statemachine proc fail")
	}

	return nil
}

func (p *OrderProcessor) DeleteFid(request *param.DeleteFidRequest) (interface{}, error) {
	rsp := &param.DeleteFidResponse{
		OrderId: request.OrderId,
		Status:  param.FAIL,
		Fids:    make(map[string]int, len(request.Fids)),
	}
	//初始化删除响应.
	p.initDeleteRsp(rsp, request)

	strategy, err := utils.GetFidDeleteStrategy(&param.GetFidDeleteStrategyRequest{
		OrderId: request.OrderId,
		Fids:    request.Fids,
	})

	if err != nil {
		utils.Log(utils.WARN, "OrderProcessor DeleteFid GetObjectDeleteStrategy fail", err.Error(), request)
		return rsp, nil
	}

	state, err := p.stateMachine.GetOrderStateInfo(request.OrderId)
	if err != nil { //从数据库加载。
		state, err = p.stateMachine.GetStateFromDB(request.OrderId)
		if err != nil {
			utils.Log(utils.WARN, "OrderProcessor DeleteFid GetObjectDeleteStrategy fail", err.Error(), request)
			return rsp, nil
		}
	}

	if state.Status != dict.TASK_CHARGE_SUC { //非正常订单，不删除。
		utils.Log(utils.WARN, "OrderProcessor DeleteFid state.Status != dict.TASK_CHARGE_SUC", "", request)
		return rsp, nil
	}

	//to do delete, 或者调整备份策略。
	deleteRequest, repRequest := p.generateRequest(request.OrderId, strategy, state)

	if len(deleteRequest.Tasks) > 0 {
		deleteRsp, err := utils.DeleteFid(deleteRequest)
		if err != nil {
			utils.Log(utils.WARN, "OrderProcessor DeleteFid DeleteOrderFid fail ", err.Error(), deleteRequest)
			return rsp, nil
		}

		p.updateOrderDeleteState(deleteRsp, state)

		p.updateFidDelRsp(deleteRsp, rsp.Fids)
	}

	if len(repRequest.Tasks) > 0 {
		repRsp, err := utils.Replicate(repRequest)
		if err != nil {
			utils.Log(utils.WARN, "OrderProcessor DeleteFid Replicate fail ", err.Error(), deleteRequest)
			return rsp, nil
		}

		p.updateOrderReps(repRequest, repRsp, state)

		p.updateFidRepRsp(repRsp, rsp.Fids)
	}

	if err := p.stateMachine.UpdateOrderStateInfo(request.OrderId, state); err != nil {
		utils.Log(utils.ERROR, "OrderProcessor DeleteFid UpdateOrderStateInfo fail ", err.Error(), deleteRequest)
		return rsp, nil
	}

	p.stateMachine.Delete(request.OrderId)

	rsp.Status = param.SUCCESS
	return rsp, nil
}

func (p *OrderProcessor) generateRequest(orderId string, strategy *param.StrategyInfo, state *dict.OrderStateInfo) (deleteRequest *param.DeleteOrderFidRequest, repRequest *param.ReplicationRequest) {
	deleteRequest = &param.DeleteOrderFidRequest{
		OrderId: orderId,
		Tasks:   make(map[string][]*param.UploadTask, 2),
	}

	repRequest = &param.ReplicationRequest{
		OrderId: orderId,
		Tasks:   make(map[string]*dict.Task, 2),
	}

	for _, task := range strategy.Tasks { //删除失败，重新调用.
		if stateTask, ok := state.Tasks[task.Fid]; ok {
			for _, rep := range task.Reps {
				stateRep, ok := stateTask.Reps[rep.Region]
				if !ok || stateRep.Status == dict.TASK_DEL_SUC || rep.Status == dict.TASK_DEL_SUC { //备份区域在订单状态中不存在，或者已经删除成功了,则过滤掉。
					if ok {
						stateRep.Status = dict.TASK_DEL_SUC
					}

					continue
				}

				//没有删除。
				if rep.MaxRep == 0 { //删除备份
					if _, ok := deleteRequest.Tasks[task.Fid]; !ok {
						deleteRequest.Tasks[task.Fid] = make([]*param.UploadTask, 0, 5)
					}

					deleteRequest.Tasks[task.Fid] = append(deleteRequest.Tasks[task.Fid], &param.UploadTask{
						Fid:    stateTask.Fid,
						Cid:    stateTask.Cid,
						Region: rep.Region,
					})

				} else if rep.MaxRep > 0 { //调整备份数

					repTask, ok := repRequest.Tasks[task.Fid]
					if !ok {
						repTask = &dict.Task{
							Fid:    stateTask.Fid,
							Cid:    stateTask.Cid,
							Region: stateTask.Region,
							Reps:   make(map[string]*dict.Rep),
						}

						repRequest.Tasks[task.Fid] = repTask
					}
					newRep := *rep
					repTask.Reps[rep.Region] = &newRep //按策略调整备份数.

					//调整订单状态分区中的最小最大备份数.
					//stateRep.MinRep = newRep.MinRep
					//stateRep.MaxRep = newRep.MaxRep
				}
			}
		}
	}
	return
}

func (p *OrderProcessor) updateOrderDeleteState(deleteRsp *param.DeleteOrderFidResponse, state *dict.OrderStateInfo) {
	for fid, tasks := range deleteRsp.Tasks {
		stateTask, ok := state.Tasks[fid]
		if !ok {
			continue
		}
		for _, task := range *tasks {
			if rep, ok := stateTask.Reps[task.Region]; ok {
				if task.Status == param.SUCCESS {
					rep.Status = dict.TASK_DEL_SUC
				}
			}
		}
	}
}

func (p *OrderProcessor) updateFidDelRsp(deleteRsp *param.DeleteOrderFidResponse, fids map[string]int) {
	for fid, tasks := range deleteRsp.Tasks {
		fids[fid] = param.SUCCESS
		for _, task := range *tasks { //一个region 没有删除成功，则删除失败。
			if task.Status == param.FAIL {
				fids[fid] = param.FAIL
				break
			}
		}
	}
}

func (p *OrderProcessor) updateOrderReps(repRequest *param.ReplicationRequest, repRsp *param.ReplicationResponse, state *dict.OrderStateInfo) {
	for _, task := range repRsp.Tasks {
		if repTask, ok := repRequest.Tasks[task.Fid]; ok {
			if stateTask, ok := state.Tasks[task.Fid]; ok {
				for region, status := range task.RegionStatus {
					if status != param.SUCCESS {
						continue
					}

					if repRep, ok := repTask.Reps[region]; ok {
						if stateRep, ok := stateTask.Reps[region]; ok {
							stateRep.MinRep = repRep.MinRep
							stateRep.MaxRep = repRep.MaxRep
						}
					}
				}
			}
		}
	}
}

func (p *OrderProcessor) updateFidRepRsp(repRsp *param.ReplicationResponse, fids map[string]int) {
	for _, task := range repRsp.Tasks { //如果一个区域执行备份失败，则订单备份失败。
		fids[task.Fid] = param.SUCCESS
		for _, status := range task.RegionStatus {
			if status == param.FAIL { //一个集群调整备份数失败，则fid执行删除失败。
				fids[task.Fid] = param.FAIL
				break
			}
		}
	}
}

//默认为成功.
func (p *OrderProcessor) initDeleteRsp(rsp *param.DeleteFidResponse, request *param.DeleteFidRequest) {
	for fid, _ := range request.Fids {
		rsp.Fids[fid] = param.SUCCESS
	}
}
