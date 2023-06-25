package processor

import (
	"bytes"
	"controller/business/config"
	"controller/business/database"
	"controller/business/dict"
	"controller/business/param"
	ctl "controller/pkg/http"
	"controller/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Order struct {
	uploadRequest   *database.UploadRequest
	downloadRequest *database.DownloadRequest
	task            *database.Task
	domainMap       map[string]*dict.DomainInfo
}

func (p *Order) Init(db *database.DataBase) {
	p.uploadRequest = new(database.UploadRequest)
	p.uploadRequest.Init(db)

	p.downloadRequest = new(database.DownloadRequest)
	p.downloadRequest.Init(db)

	p.task = new(database.Task)
	p.task.Init(db)
	//init domain

	domain := new(database.Domain)
	domain.Init(db)

	domains, err := domain.Load()
	if err != nil {
		panic(fmt.Sprintf("domian load fail, err: %v", err.Error()))
	}

	p.domainMap = make(map[string]*dict.DomainInfo, len(domains))
	for i, _ := range domains {
		p.domainMap[domains[i].Region] = &domains[i]
	}

}

func (p *Order) CreateUploadOrder(request *param.UploadTaskRequest) (interface{}, error) {
	if enough, err := p.CheckAccountBalance(&param.CheckBalanceRequest{UserId: request.UserId}); err != nil || enough == false {
		if err != nil {
			return nil, err
		}
		if enough == false {
			return nil, errors.New("insufficient fund")
		}
	}

	nodes, err := p.getNodeList(request)
	if err != nil {
		return nil, err
	}

	if err := p.saveRequest(request); err != nil {
		return nil, err
	}

	orderTask := &param.OrderTaskResponse{}
	orderTask, err = p.createOrderTask(&param.CreateTaskRequest{RequestId: request.RequestId, Type: param.UPLOAD})
	if err != nil {
		return nil, err
	}

	if err := p.saveTaskInfo(orderTask.OrderId, request); err != nil {
		return nil, err
	}

	_, err = p.createStrategy(&param.CreateStrategyRequest{RequestId: request.RequestId, OrderId: orderTask.OrderId})
	if err != nil {
		return nil, err
	}

	return &param.UploadTaskResponse{
		Status:   param.SUCCESS,
		OrderId:  orderTask.OrderId,
		NodeList: nodes,
	}, nil
}

func (p *Order) UploadFinish(request *param.UploadFinishRequest) (interface{}, error) {
	return p.uploadFinish(request)
}

func (p *Order) DownloadFinish(request *param.DownloadFinishRequest) (interface{}, error) {
	return p.downloadFinish(request)
}

func (p *Order) saveRequest(request *param.UploadTaskRequest) error {
	orgRequestInfo := p.generateUploadRequestInfo(request)
	//判断request_id是否已经存在
	count, err := p.uploadRequest.GetOrgRequestCount(request.RequestId)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("org request have exist")
	}

	return p.uploadRequest.Add(orgRequestInfo)
}

func (p *Order) saveDownloadRequest(request *param.DownloadTaskRequest) error {
	count, err := p.downloadRequest.GetDownloadRequestCount(request.RequestId)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(fmt.Sprintf("saveDownloadRequest requestId: %v hava exist", request.RequestId))
	}

	downloadRequest := p.generateDownloadRequestInfo(request)

	return p.downloadRequest.Add(downloadRequest)
}

func (p *Order) saveTaskInfo(orderId string, request *param.UploadTaskRequest) error {
	//已经存在.
	count, err := p.task.GetRequestTaskCount(request.RequestId)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(fmt.Sprintf("uploadtask saveTaskInfo request id: %v, task have exist", request.RequestId))
	}

	tasks := p.generateTaskInfo(orderId, request)

	for _, task := range tasks {
		if err := p.task.Add(task); err != nil {
			logger.Warn(" saveTaskInfo to db fail, err: %v, task:%v", err.Error(), *task)
			return err
		}
	}

	return nil
}

func (p *Order) getNodeList(request *param.UploadTaskRequest) ([]string, error) {
	if request.UploadType == param.HaveDevice {
		return request.NasList, nil
	}

	nodes, err := p.getNodeListFromRM(&param.NodeListRequst{Group: request.Group, Tag: ""})
	if err != nil {
		return nil, err
	}
	if nodes.Status != param.SUCCESS {
		return nil, errors.New("get response fail")
	}

	return nodes.Knodes, nil
}

func (p *Order) getNodeListFromRM(request *param.NodeListRequst) (*param.NodeListResponse, error) {
	domain, ok := p.domainMap[request.Group]
	if !ok {
		return nil, errors.New(fmt.Sprintf("getNodeListFromRM group: %v not exit in domian ", domain))
	}

	url := fmt.Sprintf("%s://%s/api/v0/nodelist", config.ServerCfg.Request.Protocol, domain.Url)

	queryParam := make(map[string]string, 3)
	queryParam["group"] = request.Group
	queryParam["tag"] = request.Tag

	rsp, err1 := ctl.DoRequest(http.MethodGet, url, queryParam, nil)
	if err1 != nil {
		return nil, err1
	}

	ret := &param.NodeListResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}
	if ret.Status != param.SUCCESS {
		return nil, errors.New("get nodelist from rm fail")
	}

	return ret, nil
}

func (p *Order) createOrderTask(request *param.CreateTaskRequest) (*param.OrderTaskResponse, error) {
	nameServerURL := fmt.Sprintf("%s://%s/task_tracker/v1/createTask", config.ServerCfg.Request.Protocol, config.ServerCfg.TaskTracker.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	rsp, err1 := ctl.DoRequest(http.MethodPost, nameServerURL, nil, bytes.NewReader(bt))
	if err1 != nil {
		return nil, err1
	}

	ret := &param.OrderTaskResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}

	if ret.Status != param.SUCCESS {
		return nil, errors.New("create order task fail")
	}

	return ret, nil
}

func (p *Order) createStrategy(request *param.CreateStrategyRequest) (*param.CreateStrategyResponse, error) {
	url := fmt.Sprintf("%s://%s/strategy/v1/createStrategy", config.ServerCfg.Request.Protocol, config.ServerCfg.Strategy.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	rsp, err1 := ctl.DoRequest(http.MethodPost, url, nil, bytes.NewReader(bt))
	if err1 != nil {
		return nil, err1
	}

	ret := &param.CreateStrategyResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}
	if ret.Status != param.SUCCESS {
		return nil, errors.New("create strategy fail")
	}

	return ret, nil
}

func (p *Order) CheckAccountBalance(request *param.CheckBalanceRequest) (bool, error) {
	url := fmt.Sprintf("%s://%s/account/v1/checkBalance", config.ServerCfg.Request.Protocol, config.ServerCfg.Account.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return false, err
	}

	rsp, err1 := ctl.DoRequest(http.MethodPost, url, nil, bytes.NewReader(bt))
	if err1 != nil {
		return false, err1
	}

	ret := &param.CheckBalanceResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return false, err
	}
	if ret.Status == param.SUCCESS {
		return ret.Enough, nil
	} else {
		return false, errors.New("get account response fail")
	}
}

func (p *Order) uploadFinish(request *param.UploadFinishRequest) (interface{}, error) {
	url := fmt.Sprintf("%s://%s/task_tracker/v1/uploadFinish", config.ServerCfg.Request.Protocol, config.ServerCfg.TaskTracker.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return false, err
	}

	rsp, err1 := ctl.DoRequest(http.MethodPost, url, nil, bytes.NewReader(bt))
	if err1 != nil {
		return false, err1
	}

	ret := &param.UploadFinishResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return false, err
	}

	return ret, nil
}

func (p *Order) generateUploadRequestInfo(req *param.UploadTaskRequest) *dict.UploadRequestInfo {
	if req == nil {
		return nil
	}

	uploadRequestInfo := &dict.UploadRequestInfo{
		RequestId:  req.RequestId,
		UserId:     req.UserId,
		UploadType: req.UploadType,
		Tasks:      make([]*dict.UploadTask, 0, len(req.Tasks)),
		Group:      req.Group,
		NasList:    req.NasList,
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	}

	for _, task := range req.Tasks {
		uploadRequestInfo.Tasks = append(uploadRequestInfo.Tasks, task)
	}

	return uploadRequestInfo
}

func (p *Order) generateDownloadRequestInfo(req *param.DownloadTaskRequest) *dict.DownloadRequestInfo {
	if req == nil {
		return nil
	}

	downloadRequestInfo := &dict.DownloadRequestInfo{
		RequestId:   req.RequestId,
		UserId:      req.UserId,
		DownlodType: req.DownloadType,
		Tasks:       make([]*dict.DownloadTask, 0, len(req.Tasks)),
		CreateTime:  time.Now().UnixMilli(),
		UpdateTime:  time.Now().UnixMilli(),
	}

	for _, task := range req.Tasks {
		downloadRequestInfo.Tasks = append(downloadRequestInfo.Tasks, task)
	}

	return downloadRequestInfo
}

func (p *Order) generateTaskInfo(orderId string, req *param.UploadTaskRequest) []*dict.TaskInfo {
	if req == nil {
		return nil
	}

	tasks := make([]*dict.TaskInfo, 0, len(req.Tasks))
	for _, task := range req.Tasks {
		tasks = append(tasks, &dict.TaskInfo{
			Fid:        task.Fid,
			Cid:        "",
			OrderId:    orderId,
			RequestId:  req.RequestId,
			Rep:        task.Rep,
			VirtualRep: 0,
			RepMin:     task.RepMin,
			RepMax:     task.RepMax,
			Encryption: task.Encryption,
			Expire:     task.Expire,
			Level:      task.Level,
			Size:       task.Size,
			Name:       task.Name,
			Ajust:      task.Ajust,
			Desc:       "",
			Status:     0,
			CreateTime: time.Now().UnixMilli(),
			UpdateTime: time.Now().UnixMilli(),
		})
	}

	return tasks
}

func addDomainRecord(domain *database.Domain) {
	err := domain.Add(&dict.DomainInfo{
		Id:         1,
		Region:     "chengdu",
		Url:        "192.168.2.35:9094",
		Status:     dict.SUCCESS,
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	})

	if err != nil {
		logger.Errorf(" add damain record fail , err: %v", err.Error())
	}
}

func (p *Order) CreateDownloadOrder(request *param.DownloadTaskRequest) (interface{}, error) {
	if enough, err := p.CheckAccountBalance(&param.CheckBalanceRequest{UserId: request.UserId}); err != nil || enough == false {
		if err != nil {
			return nil, err
		}
		if enough == false {
			return nil, errors.New("insufficient fund")
		}
	}

	cids := p.convertCids(request.Tasks)

	nodes, err := p.getDownloadNodesFromRM(request.Group, cids)
	if err != nil {
		return nil, err
	}

	if err := p.saveDownloadRequest(request); err != nil {
		return nil, err
	}

	orderTask := &param.OrderTaskResponse{}
	orderTask, err = p.createOrderTask(&param.CreateTaskRequest{RequestId: request.RequestId, Type: param.DOWNLOAD})
	if err != nil {
		return nil, err
	}

	return &param.DownloadTaskResponse{
		Status:  param.SUCCESS,
		OrderId: orderTask.OrderId,
		Nodes:   nodes.Data,
	}, nil
}

func (p *Order) convertCids(tasks []*dict.DownloadTask) string {
	var builder strings.Builder

	nLen := len(tasks) - 1
	if nLen < 0 {
		return ""
	}

	for i := 0; i < nLen; i++ {
		builder.WriteString(tasks[i].Cid)
		builder.WriteString(",")
	}
	builder.WriteString(tasks[nLen].Cid)

	return builder.String()
}

func (p *Order) getDownloadNodesFromRM(group, cids string) (*param.GetDownloadNodeResponse, error) {
	domain, ok := p.domainMap[group]
	if !ok {
		return nil, errors.New(fmt.Sprintf("getNodeListFromRM group: %v not exit in domian ", domain))
	}

	url := fmt.Sprintf("%s://%s/api/v0/nodes/%s", config.ServerCfg.Request.Protocol, domain.Url, cids)
	rsp, err1 := ctl.DoRequest(http.MethodGet, url, nil, nil)
	if err1 != nil {
		return nil, err1
	}

	ret := &param.GetDownloadNodeResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}
	if ret.Status != param.SUCCESS {
		return nil, errors.New("get download nodelist from rm fail")
	}

	/*bt, _ := json.Marshal(ret)
	logger.Infof("helo getDownloadNodesFromRM, url:", url, "rsp:", string(bt))
	*/
	return ret, nil
}

func (p *Order) downloadFinish(request *param.DownloadFinishRequest) (interface{}, error) {
	url := fmt.Sprintf("%s://%s/task_tracker/v1/downloadFinish", config.ServerCfg.Request.Protocol, config.ServerCfg.TaskTracker.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return false, err
	}

	rsp, err1 := ctl.DoRequest(http.MethodPost, url, nil, bytes.NewReader(bt))
	if err1 != nil {
		return false, err1
	}

	ret := &param.DownloadFinishResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return false, err
	}

	return ret, nil
}
