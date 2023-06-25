package processor

import (
	"controller/task_tracker/database"
	"controller/task_tracker/dict"
	"controller/task_tracker/event"
	"controller/task_tracker/param"
	"controller/task_tracker/statemachine"
	"errors"
	"fmt"
	"time"
)

type OrderProcessor struct {
	order           *database.Order
	uploadRequest   *database.UploadRequest
	downloadRequest *database.DownloadRequest
	stateMachine    *statemachine.StateMachine
}

func (p *OrderProcessor) Init(db *database.DataBase, machine *statemachine.StateMachine) {
	p.order = new(database.Order)
	p.order.Init(db)

	p.uploadRequest = new(database.UploadRequest)
	p.uploadRequest.Init(db)

	p.downloadRequest = new(database.DownloadRequest)
	p.downloadRequest.Init(db)

	p.stateMachine = machine
}

func (p *OrderProcessor) CreateOrderTask(request *param.CreateTaskRequest) (interface{}, error) {
	if request.Type == param.UPLOAD {
		return p.createOrderUploadTask(request)
	} else if request.Type == param.DOWNLOAD {
		return p.createOrderDownloadTask(request)
	}

	return nil, errors.New(fmt.Sprintf("unknown order type: %d ", request.Type))
}

func (p *OrderProcessor) createOrderUploadTask(request *param.CreateTaskRequest) (interface{}, error) {
	uploadRequest, err := p.uploadRequest.GetOrgRequest(request.RequestId)
	if err != nil {
		return nil, err
	}

	//判断订单是否创建完成
	if orderInfo, err := p.stateMachine.GetOrderByRequestId(request.RequestId); err == nil {
		return param.UploadTaskResponse{
			Status:  param.SUCCESS,
			OrderId: orderInfo.OrderId,
		}, nil
	}

	orderInfo := p.generateOrder(request)

	event, err1 := p.generateUploadOrderEvent(orderInfo.OrderId, uploadRequest)
	if err1 != nil {
		return nil, err1
	}

	if err := p.stateMachine.Send(orderInfo.OrderId, event); err != nil {
		return nil, err
	}

	status := <-event.Ret

	return param.UploadTaskResponse{
		Status:  status,
		OrderId: orderInfo.OrderId,
	}, nil
}

func (p *OrderProcessor) createOrderDownloadTask(request *param.CreateTaskRequest) (interface{}, error) {
	downloadRequest, err := p.downloadRequest.GetDownloadRequst(request.RequestId)
	if err != nil {
		return nil, err
	}

	//判断订单是否创建完成
	if orderInfo, err := p.stateMachine.GetOrderByRequestId(request.RequestId); err == nil {
		return param.DownloadTaskResponse{
			Status:  param.SUCCESS,
			OrderId: orderInfo.OrderId,
		}, nil
	}

	orderInfo := p.generateOrder(request)

	event, err1 := p.generateDownloadOrderEvent(orderInfo.OrderId, downloadRequest)
	if err1 != nil {
		return nil, err1
	}

	if err := p.stateMachine.Send(orderInfo.OrderId, event); err != nil {
		return nil, err
	}

	status := <-event.Ret

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
		} else {
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

/*func (p *OrderProcessor) UploadFinish(request *param.UploadFinishRequest) (interface{}, error) {
	status := param.FAIL

	uploadType := 0
	if request.Status == param.SUCCESS { //等待一定时间, 如果都不成功，则删除.
		for i := 0; i < 15; i++ {
			time.Sleep(10 * time.Second)
			state, err := p.stateMachine.GetOrderStateInfo(request.OrderId)
			if err != nil {
				return nil, err
			}

			uploadType = p.judageFileUploadFinish(request, state)
			if uploadType == UPLOAD_FINISH {
				break
			}
		}
	}

	if uploadType == UPLOAD_FINISH {
		status = param.SUCCESS
	}

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
		Status:  status,
	}, nil
}
*/

func (p *OrderProcessor) judageFileUploadFinish(request *param.UploadFinishRequest, state *dict.OrderStateInfo) int {
	if len(request.Tasks) == len(state.Tasks) {
		for _, value := range request.Tasks {
			if task, ok := state.Tasks[value.Fid]; !ok { //不存在fid
				//logger.Infof("judageFileUploadFinish not finish request.Fid: %v not exit in order sate info", value.Fid)
				return FILE_NOT_EXIST
			} else {
				if task.Cid != value.Cid { //存在fid ,但cid不相等
					//logger.Infof("judageFileUploadFinish not finish ,task.Cid: %v  != value.Cid: %v", task.Cid, value.Cid)
					return CID_NOT_EQUAL
				}
			}
		}
	} else { //长度不相等
		//logger.Infof("judageFileUploadFinish not finish  len(request.Tasks) %v != len(state.Tasks): %v", len(request.Tasks), len(state.Tasks))
		return FILE_NUM_NOT_EQUAL
	}

	return UPLOAD_FINISH
}

func (p *OrderProcessor) generateOrder(request *param.CreateTaskRequest) *dict.OrderInfo {
	return &dict.OrderInfo{
		OrderId:    CreateOrderId(),
		RequestId:  request.RequestId,
		OrderType:  request.Type,
		Status:     dict.TASK_INIT,
		Desc:       "",
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli()}
}

func (p *OrderProcessor) generateUploadOrderEvent(orderId string, uploadInfo *dict.UploadRequestInfo) (*event.Event, error) {
	createOrderEvent := &event.CreateOrderEvent{RequestId: uploadInfo.RequestId, OrderType: param.UPLOAD}

	createOrderEvent.Fids = make([]string, 0, len(uploadInfo.Tasks))
	for _, task := range uploadInfo.Tasks {
		createOrderEvent.Fids = append(createOrderEvent.Fids, task.Fid)
	}

	return &event.Event{
		Type:    event.CREATE_ORDER,
		OrderId: orderId,
		Ret:     make(chan int),
		Data:    createOrderEvent,
	}, nil
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
	if request.Status == param.FAIL {
		return &param.DownloadFinishResponse{
			OrderId: request.OrderId,
			Status:  param.SUCCESS,
		}, nil
	}

	//
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
