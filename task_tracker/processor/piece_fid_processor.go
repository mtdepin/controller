package processor

import (
	"bytes"
	"controller/api"
	"controller/pkg/distributionlock"
	ctl "controller/pkg/http"
	"controller/pkg/logger"
	"controller/pkg/newcache"
	"controller/task_tracker/config"
	"controller/task_tracker/database"
	"controller/task_tracker/dict"
	"controller/task_tracker/event"
	e "controller/task_tracker/event"
	"controller/task_tracker/param"
	"controller/task_tracker/statemachine"
	"controller/task_tracker/utils"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"time"
)

type PieceFidProcessor struct {
	order           *database.Order
	uploadRequest   *database.UploadRequest
	downloadRequest *database.DownloadRequest
	fidReplicate    *database.FidReplication
	stateMachine    *statemachine.StateMachine
	cache           newcache.Cache
	mutex           *distributionlock.MutexLock
	fidEventChan    chan *e.FidEvent
	prefix          string
	nThread         int
}

func (p *PieceFidProcessor) Init(db *database.DataBase, machine *statemachine.StateMachine, fidReplicate *database.FidReplication) {
	p.order = new(database.Order)
	p.order.Init(db)

	p.uploadRequest = new(database.UploadRequest)
	p.uploadRequest.Init(db)

	p.downloadRequest = new(database.DownloadRequest)
	p.downloadRequest.Init(db)

	p.fidReplicate = fidReplicate
	p.stateMachine = machine
	p.prefix = config.ServerCfg.Lock.Prefix

	p.prefix = config.ServerCfg.Lock.Prefix

	p.nThread = FID_EVENT_SIZE

	p.fidEventChan = make(chan *e.FidEvent, p.nThread)
	for i := 0; i < p.nThread; i++ {
		go p.FidHandler()
	}

	p.initRedis()
}

func (p *PieceFidProcessor) initRedis() {
	fidAddrs := strings.Split(config.ServerCfg.Redis.FidAddress, ",") //注意两个redis cluster
	lockAddrs := strings.Split(config.ServerCfg.Redis.LockAddress, ",")
	fidPassword := config.ServerCfg.Redis.FidPassword
	lockPassword := config.ServerCfg.Redis.LockPassword

	p.cache = newcache.NewRedisClient(fidAddrs, fidPassword)
	p.mutex = distributionlock.NewMutexLock(lockAddrs, lockPassword)
	logger.Info("init redis cluster success")
}

func (p *PieceFidProcessor) UploadPieceFid(request *api.UploadPieceFidRequest) (interface{}, error) {
	pieceFids, err := p.getPieceFidInfo(request)
	if err != nil {
		return nil, err
	}

	//to do add task.
	tasks := p.getTasks(pieceFids)
	p.stateMachine.AddPieceFid(request.OrderId, tasks)

	finish, err := p.stateMachine.TaskInitFinish(request.OrderId)
	if err != nil {
		return nil, err
	}

	if finish {
		if _, err := p.createStrategy(&param.CreateStrategyRequest{OrderId: request.OrderId, Region: request.Group, Ext: &param.Extend{Ctx: request.Ext.Ctx}}); err != nil {
			return nil, err
		}
	}

	//重复fid ,发送重复事件。

	//判断重复fid
	index := 0
	repFids := make(map[string]string)
	for idx, pieceFid := range pieceFids {
		if pieceFid.Cid != "" {
			repFids[pieceFid.Fid] = pieceFid.Cid
			index = idx
		}
	}

	if finish && len(repFids) == len(request.Pieces) { //如果有非重复文件，则一定会上传， 触发备份, 如果都重复了，则不会上传了， 触发一次备份。
		event, err := p.generateCallbackUploadEvent(request.OrderId, pieceFids[index])
		if err != nil {
			return nil, err
		}

		if err := p.stateMachine.Send(event.OrderId, event); err != nil {
			return nil, err
		}
		status := <-event.Ret
		if status != param.SUCCESS {
			bt, _ := json.Marshal(pieceFids[index])
			logger.Errorf("UploadPieceFid, stateMachine send repeat file upload success fail: %v, orderId: %v, fidInfo: %v", status, request.OrderId, string(bt))
		}
	}

	//to do proc repeate file

	return &api.UploadPieceFidResponse{
		OrderId: request.OrderId,
		RepFids: repFids,
		Status:  param.SUCCESS,
	}, nil
}

func (p *PieceFidProcessor) getTasks(pieceFids []*event.FidInfo) []*dict.Task {
	tasks := make([]*dict.Task, 0, len(pieceFids))
	for _, pieceFid := range pieceFids {
		task := &dict.Task{Fid: pieceFid.Fid, Cid: pieceFid.Cid, Status: pieceFid.Status, Repeate: pieceFid.Repeate, Origins: pieceFid.Origins, Region: pieceFid.Region}
		tasks = append(tasks, task)
	}

	return tasks
}

func (p *PieceFidProcessor) getPieceFidInfo(request *api.UploadPieceFidRequest) ([]*event.FidInfo, error) {
	nLen := len(request.Pieces)
	if nLen == 0 {
		return nil, errors.New("upload tasks is empty")
	}

	fidInfos := make([]*event.FidInfo, 0, nLen)
	rets := make([]chan *e.FidRet, 0, nLen)

	count := 0
	var err error

	for _, task := range request.Pieces {
		count++
		ret := make(chan *e.FidRet)

		p.fidEventChan <- &e.FidEvent{
			Fid:   task.Fid,
			Group: request.Group,
			Ret:   ret,
		}

		rets = append(rets, ret)

		if count >= p.nThread-1 {
			for _, ret := range rets {
				fidRet := <-ret
				if fidRet.Err != nil {
					err = fidRet.Err
				}
				fidInfos = append(fidInfos, fidRet.FidInfo)
			}
			//重置
			count = 0
			rets = rets[0:0]

			if err != nil {
				return nil, err
			}
		}
	}

	for _, ret := range rets {
		fidRet := <-ret
		if fidRet.Err != nil {
			err = fidRet.Err
		}
		fidInfos = append(fidInfos, fidRet.FidInfo)
	}

	if err != nil {
		return nil, err
	}
	return fidInfos, nil
}

func (p *PieceFidProcessor) FidHandler() {
	for {
		event := <-p.fidEventChan
		fidInfo, err := p.getFidInfo(event)
		event.Ret <- &e.FidRet{FidInfo: fidInfo, Err: err}
	}
}

func (p *PieceFidProcessor) getFidInfo(event *e.FidEvent) (*e.FidInfo, error) {
	ret := &e.FidInfo{
		Fid:     event.Fid,
		Cid:     "",
		Repeate: 0,
		Origins: "",
		Status:  dict.TASK_INIT,
	}

	fidInfo, err := p.fidReplicate.Search(event.Fid)
	if err != nil { //不存在，直接返回.
		return ret, nil
	}

	if fidInfo.Cid == "" || fidInfo.Status != dict.TASK_REP_SUC { //修改:  备份成功作为重复文件。
		return ret, nil
	}

	key := p.prefix + event.Fid
	if err := p.mutex.Lock(key); err != nil {
		utils.Log(utils.ERROR, "PieceFidProcessor generateUploadOrderEvent umutex.Lock(key): %v, fail:%v ", key, err.Error())
		return ret, err
	}
	defer p.mutex.UnLock(key)

	ret.Cid = fidInfo.Cid
	ret.Status = dict.TASK_UPLOAD_SUC
	ret.Repeate = dict.REPEATE
	ret.Origins = fidInfo.Origins
	ret.Region = fidInfo.Region

	if err := p.updateFidInfo(event.Fid, event.Group, fidInfo, 1); err != nil {
		utils.Log(utils.ERROR, "PieceFidProcessor generateUploadOrderEvent update updateFidInfo fail", err.Error(), fidInfo)
		return ret, err
	}

	return ret, nil
}

func (p *PieceFidProcessor) updateFidInfo(fid, region string, fidInfo *dict.FidInfo, used int) error {
	//fidInfo.Used = used
	fidInfo.Region = region
	fidInfo.UpdateTime = time.Now().UnixMilli()
	if reps, ok := fidInfo.Reps[region]; ok {
		for _, rep := range reps { //设置任意订单。
			rep.Used = used
			break
		}
	} else { //不存在，设置任意区域，任意订单的cid 被占用.
		for _, reps := range fidInfo.Reps {
			for _, rep := range reps {
				rep.Used = used
				break
			}
			break
		}
	}

	bt, err := json.Marshal(fidInfo)
	if err != nil {
		logger.Error(fmt.Sprintf("updateFidInfo, json marshal fidInfo fail: %v, fidInfo: %v", err.Error(), fidInfo))
		return err
	}

	if err := p.cache.Set(fid, string(bt), 0); err != nil {
		return err
	}

	if err := p.fidReplicate.Update(fid, bson.M{"$set": bson.M{"reps": fidInfo.Reps}}); err != nil {
		return err
	}
	return nil
}

func (p *PieceFidProcessor) createStrategy(request *param.CreateStrategyRequest) (*param.CreateStrategyResponse, error) {
	url := fmt.Sprintf("%s://%s/strategy/v1/createStrategy", config.ServerCfg.Request.Protocol, config.ServerCfg.Strategy.Url)

	ctx := request.Ext.Ctx
	request.Ext = nil
	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	rsp, err := ctl.DoRequest(ctx, http.MethodPost, url, nil, bytes.NewReader(bt))
	if err != nil {
		return nil, err
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

func (p *PieceFidProcessor) generateCallbackUploadEvent(orderId string, fidInfo *event.FidInfo) (*e.Event, error) {
	callbackUploadEvent := &e.CallbackUploadEvent{
		OrderId: orderId,
		Fid:     fidInfo.Fid,
		Cid:     fidInfo.Cid,
		Region:  fidInfo.Region,
		Origins: fidInfo.Origins,
		Status:  param.SUCCESS,
	}

	return &e.Event{
		Type:    e.CALLBACK_UPLOAD,
		OrderId: orderId,
		Ret:     make(chan int),
		Data:    callbackUploadEvent,
	}, nil
}
