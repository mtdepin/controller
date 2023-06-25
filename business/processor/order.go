package processor

import (
	"bytes"
	"context"
	"controller/api"
	"controller/business/config"
	"controller/business/database"
	"controller/business/dict"
	"controller/business/param"
	"controller/business/utils"
	ctl "controller/pkg/http"
	"controller/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Order struct {
	uploadRequest   *database.UploadRequest
	downloadRequest *database.DownloadRequest
	task            *database.Task
	fidReplicate    *database.FidReplication
	domainMap       map[string]*dict.DomainInfo
}

func (p *Order) Init(db *database.DataBase) {
	p.uploadRequest = new(database.UploadRequest)
	p.uploadRequest.Init(db)

	p.downloadRequest = new(database.DownloadRequest)
	p.downloadRequest.Init(db)

	p.fidReplicate = new(database.FidReplication)
	p.fidReplicate.Init(db)

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

func (p *Order) CreateUploadOrder(request *api.UploadTaskRequest) (interface{}, error) {
	if enough, err := p.CheckAccountBalance(&api.CheckBalanceRequest{UserId: request.UserId, Ext: request.Ext}); err != nil || enough == false {
		if err != nil {
			return nil, err
		}
		if enough == false {
			return nil, errors.New("insufficient fund")
		}
	}

	region, err := p.selectUploadRegion(request)
	if err != nil {
		return nil, err
	}

	request.Group = region

	//set region.
	nodes, err := p.getNodeList(request)
	if err != nil {
		return nil, err
	}

	if len(nodes) == 0 { //没有上传节点，创建任务失败.
		return nil, errors.New("upload node not exist")
	}

	if err := p.saveRequest(request); err != nil {
		return nil, err
	}

	orderTask, err := p.createOrderTask(&api.CreateTaskRequest{RequestId: request.RequestId, Type: param.UPLOAD, Ext: request.Ext})
	if err != nil {
		return nil, err
	}

	return &api.UploadTaskResponse{
		Status:   param.SUCCESS,
		OrderId:  orderTask.OrderId,
		NodeList: nodes,
		Group:    region,
	}, nil
}

func (p *Order) UploadPieceFid(request *api.UploadPieceFidRequest) (interface{}, error) {
	if err := p.saveTaskInfo(request.OrderId, request); err != nil {
		return nil, err
	}

	rsp, err := p.uploadPieceFid(request)
	if err != nil {
		return nil, err
	}

	return &api.UploadPieceFidResponse{
		OrderId: request.OrderId,
		RepFids: rsp.RepFids,
		Status:  param.SUCCESS,
	}, nil
}

func (p *Order) UploadFinish(request *param.UploadFinishRequest) (interface{}, error) {
	return p.uploadFinish(request)
}

func (p *Order) DownloadFinish(request *param.DownloadFinishRequest) (interface{}, error) {
	return p.downloadFinish(request)
}

func (p *Order) saveRequest(request *api.UploadTaskRequest) error {
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

func (p *Order) saveTaskInfo(orderId string, request *api.UploadPieceFidRequest) error {
	if len(request.Pieces) == 0 {
		return errors.New("saveTaskInfo fail, len(request.Pieces) == 0 : order_id:" + orderId)
	}

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
			bt, _ := json.Marshal(task)
			logger.Warnf("saveTaskInfo to db fail, orderId: %v, err: %v, task:%v", orderId, err.Error(), string(bt))
			return err
		}
	}

	return nil
}

func (p *Order) getNodeList(request *api.UploadTaskRequest) ([]*api.Node, error) {
	if request.UploadType == param.HaveDevice {
		return request.NasList, nil
	}

	nodes, err := p.getNodeListFromRM(&api.GetKNodesRequest{Group: request.Group, NodeNum: request.NodeNum, Ext: request.Ext})
	if err != nil {
		return nil, err
	}
	if nodes.Status != param.SUCCESS {
		return nil, errors.New("get response fail")
	}

	return nodes.Knodes, nil
}

func (p *Order) getNodeListFromRM(request *api.GetKNodesRequest) (*api.NodeListResponse, error) {
	domain, ok := p.domainMap[request.Group]
	if !ok {
		bt, _ := json.Marshal(p.domainMap)
		return nil, errors.New(fmt.Sprintf("getUploadNodeListFromRM group: %v not exit in domian: %v", request.Group, string(bt)))
	}

	url := fmt.Sprintf("%s://%s/api/v0/nodelist", config.ServerCfg.Request.Protocol, domain.Url)

	queryParam := make(map[string]string, 3)
	queryParam["group"] = request.Group
	queryParam["tag"] = ""
	queryParam["node_num"] = strconv.Itoa(request.NodeNum)

	rsp, err := ctl.DoRequest(request.Ext.Ctx, http.MethodGet, url, queryParam, nil)
	if err != nil {
		return nil, err
	}

	ret := &api.NodeListResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}
	if ret.Status != param.SUCCESS {
		return nil, errors.New("get nodelist from rm fail")
	}

	logger.Infof("helo getUploadNodelist: region: %v, request.NodeNum: %v, url: %v, response: %v", request.Group, request.NodeNum, url, string(rsp))

	if len(ret.Knodes) == 0 {
		ret.Status = param.FAIL //如果节点为空，设置返回失败。
		logger.Warn("get node from rm is empty, url: ", url)
	}

	return ret, nil
}

func (p *Order) createOrderTask(request *api.CreateTaskRequest) (*api.CreateTaskResponse, error) {
	nameServerURL := fmt.Sprintf("%s://%s/task_tracker/v1/createTask", config.ServerCfg.Request.Protocol, config.ServerCfg.TaskTracker.Url)

	ctx := request.Ext.Ctx
	request.Ext = nil
	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	rsp, err := ctl.DoRequest(ctx, http.MethodPost, nameServerURL, nil, bytes.NewReader(bt))
	if err != nil {
		return nil, err
	}

	ret := &api.CreateTaskResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}

	if ret.Status != param.SUCCESS {
		return nil, errors.New("create order task fail")
	}

	return ret, nil
}

func (p *Order) uploadPieceFid(request *api.UploadPieceFidRequest) (*api.UploadPieceFidResponse, error) {
	nameServerURL := fmt.Sprintf("%s://%s/task_tracker/v1/uploadPieceFid", config.ServerCfg.Request.Protocol, config.ServerCfg.TaskTracker.Url)

	ctx := request.Ext.Ctx
	request.Ext = nil
	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	rsp, err1 := ctl.DoRequest(ctx, http.MethodPost, nameServerURL, nil, bytes.NewReader(bt))
	if err1 != nil {
		return nil, err1
	}

	ret := &api.UploadPieceFidResponse{}
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

	rsp, err1 := ctl.DoRequest(request.Ext.Ctx, http.MethodPost, url, nil, bytes.NewReader(bt))
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

func (p *Order) CheckAccountBalance(request *api.CheckBalanceRequest) (bool, error) {
	url := fmt.Sprintf("%s://%s/account/v1/checkBalance", config.ServerCfg.Request.Protocol, config.ServerCfg.Account.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return false, err
	}

	rsp, err1 := ctl.DoRequest(request.Ext.Ctx, http.MethodPost, url, nil, bytes.NewReader(bt))
	if err1 != nil {
		return false, err1
	}

	ret := &api.CheckBalanceResponse{}
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

	rsp, err1 := ctl.DoRequest(request.Ext.Ctx, http.MethodPost, url, nil, bytes.NewReader(bt))
	if err1 != nil {
		return false, err1
	}

	ret := &param.UploadFinishResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return false, err
	}

	return ret, nil
}

func (p *Order) generateUploadRequestInfo(req *api.UploadTaskRequest) *dict.UploadRequestInfo {
	if req == nil {
		return nil
	}

	uploadRequestInfo := &dict.UploadRequestInfo{
		RequestId:  req.RequestId,
		UserId:     req.UserId,
		UploadType: req.UploadType,
		PieceNum:   req.PieceNum,
		Group:      req.Group,
		Name:       req.Name,
		Size:       req.Size,
		RemoteIp:   req.RemoteIp,
		NasList:    make([]*dict.Node, 0, len(req.NasList)),
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	}
	for _, node := range req.NasList {
		uploadRequestInfo.NasList = append(uploadRequestInfo.NasList, &dict.Node{
			Address: node.Address,
			Weight:  node.Weight,
			RtcId:   node.RtcId,
		})
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

func (p *Order) generateTaskInfo(orderId string, req *api.UploadPieceFidRequest) []*dict.TaskInfo {
	if req == nil {
		return nil
	}

	tasks := make([]*dict.TaskInfo, 0, len(req.Pieces))
	for _, task := range req.Pieces {
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

func (p *Order) CreateDownloadOrder(request *param.DownloadTaskRequest) (interface{}, error) {
	if enough, err := p.CheckAccountBalance(&api.CheckBalanceRequest{UserId: request.UserId, Ext: &api.Extend{Ctx: request.Ext.Ctx}}); err != nil || enough == false {
		if err != nil {
			return nil, err
		}
		if enough == false {
			return nil, errors.New("insufficient fund")
		}
	}

	cids, mCidFid := p.parseTask(request.Tasks)

	nodes, err := p.getDownloadNodes(request.Ext.Ctx, request.Group, cids, mCidFid)
	if err != nil {
		return nil, err
	}

	if err := p.saveDownloadRequest(request); err != nil {
		return nil, err
	}

	orderTask, err := p.createOrderTask(&api.CreateTaskRequest{RequestId: request.RequestId, Type: param.DOWNLOAD, Ext: &api.Extend{Ctx: request.Ext.Ctx}})
	if err != nil {
		return nil, err
	}

	return &param.DownloadTaskResponse{
		Status:  param.SUCCESS,
		OrderId: orderTask.OrderId,
		Nodes:   nodes.Data,
	}, nil
}

func (p *Order) parseTask(tasks []*dict.DownloadTask) (string, map[string]string) {
	var builder strings.Builder
	nLen := len(tasks) - 1
	mCidFid := make(map[string]string, nLen)

	if nLen < 0 {
		return "", mCidFid
	}

	for i := 0; i < nLen; i++ {
		builder.WriteString(tasks[i].Cid)
		mCidFid[tasks[i].Cid] = tasks[i].Fid
		builder.WriteString(",")
	}
	builder.WriteString(tasks[nLen].Cid)
	mCidFid[tasks[nLen].Cid] = tasks[nLen].Fid

	return builder.String(), mCidFid
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

func (p *Order) getDownloadNodes(ctx context.Context, group, cids string, mCidFid map[string]string) (*param.GetDownloadNodeResponse, error) {
	downloadNodes, err := p.getDownloadNodesFromRM(ctx, group, cids)
	if err != nil {
		return nil, err
	}

	for cid, nodes := range downloadNodes.Data {
		notExist := true
		for _, node := range nodes { //权重都为0 则不存在。
			if node.Weight != 0 {
				notExist = false //存在.
				break
			}
		}

		//to do seach cid and upload region
		if notExist {
			fidInfo, err := p.fidReplicate.Search(mCidFid[cid])
			if err != nil {
				return nil, errors.New("err: " + err.Error() + "cid: " + cid)
			}

			cidAddr, err := p.getDownloadNodesFromRM(ctx, fidInfo.Region, cid)
			if err != nil {
				return nil, errors.New("err: " + err.Error() + "cid: " + cid)
			}

			downloadNodes.Data[cid] = cidAddr.Data[cid]
		}
	}
	//to do print
	//bt, _ := json.Marshal(downloadNodes)
	//logger.Infof("helo_download_nodelist: %v", string(bt))

	return downloadNodes, nil
}

func (p *Order) getDownloadNodesFromRM(ctx context.Context, group, cids string) (*param.GetDownloadNodeResponse, error) {
	domain, ok := p.domainMap[group]
	if !ok {
		bt, _ := json.Marshal(p.domainMap)
		return nil, errors.New(fmt.Sprintf("helo_getDownloadNodesFromRM group: %v not exit in domian: %v", group, string(bt)))
	}

	url := fmt.Sprintf("%s://%s/api/v0/nodes/%s", config.ServerCfg.Request.Protocol, domain.Url, cids)
	rsp, err1 := ctl.DoRequest(ctx, http.MethodGet, url, nil, nil)
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

	return ret, nil
}

func (p *Order) downloadFinish(request *param.DownloadFinishRequest) (interface{}, error) {
	url := fmt.Sprintf("%s://%s/task_tracker/v1/downloadFinish", config.ServerCfg.Request.Protocol, config.ServerCfg.TaskTracker.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return false, err
	}

	rsp, err1 := ctl.DoRequest(request.Ext.Ctx, http.MethodPost, url, nil, bytes.NewReader(bt))
	if err1 != nil {
		return false, err1
	}

	ret := &param.DownloadFinishResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return false, err
	}

	return ret, nil
}

func (p *Order) DeleteFid(request *param.DeleteFidRequest) (interface{}, error) {
	url := fmt.Sprintf("%s://%s/task_tracker/v1/deleteFid", config.ServerCfg.Request.Protocol, config.ServerCfg.TaskTracker.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return false, err
	}

	rsp, err := ctl.DoRequest(request.Ext.Ctx, http.MethodPost, url, nil, bytes.NewReader(bt))
	if err != nil {
		return false, err
	}

	ret := &param.DeleteFidResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return false, err
	}

	return ret, nil
}

func (p *Order) selectUploadRegion(request *api.UploadTaskRequest) (string, error) {
	rsp, err := utils.GetRmSpaceInfo(request.Ext.Ctx, request.Group)
	if err != nil {
		return "", err
	}

	if rsp.Region == nil {
		return "", errors.New("selectUploadRegion GetRmSpaceInfo  region is nul")
	}

	if rsp.Region.ValidStorage > dict.RM_LFTESPACE_THRESHOLD {
		return request.Group, nil
	}

	if _, ok := p.domainMap[config.ServerCfg.SuperCluster.Region]; !ok {
		return "", errors.New("super region not exist")
	}

	return config.ServerCfg.SuperCluster.Region, nil
}
