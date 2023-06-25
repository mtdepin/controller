package utils

import (
	"bytes"
	"context"
	ctl "controller/pkg/http"
	"controller/pkg/logger"
	"controller/task_tracker/config"
	"controller/task_tracker/dict"
	"controller/task_tracker/param"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	INFO = iota + 1
	WARN
	ERROR
)

func Replicate(request *param.ReplicationRequest) (*param.ReplicationResponse, error) {
	nameServerURL := fmt.Sprintf("%s://%s/scheduler/v1/replicate", config.ServerCfg.Request.Protocol, config.ServerCfg.Scheduler.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	rsp, err1 := ctl.DoRequest(context.Background(), http.MethodPost, nameServerURL, nil, bytes.NewReader(bt))
	if err1 != nil {
		return nil, err1
	}

	ret := &param.ReplicationResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}

	if ret.Status != param.SUCCESS {
		return nil, errors.New("replicate fail")
	}

	return ret, nil
}

func Delete(request *param.DeleteOrderRequest) (*param.DeleteOrderResponse, error) {
	nameServerURL := fmt.Sprintf("%s://%s/scheduler/v1/deleteOrder", config.ServerCfg.Request.Protocol, config.ServerCfg.Scheduler.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	rsp, err1 := ctl.DoRequest(request.Ext.Ctx, http.MethodPost, nameServerURL, nil, bytes.NewReader(bt))
	if err1 != nil {
		return nil, err1
	}

	ret := &param.DeleteOrderResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}

	if ret.Status != param.SUCCESS {
		return nil, errors.New("delete order upload task  fail")
	}

	return ret, nil
}

func Charge(request *param.ChargeRequest) (*param.ChargeResponse, error) {
	nameServerURL := fmt.Sprintf("%s://%s/scheduler/v1/charge", config.ServerCfg.Request.Protocol, config.ServerCfg.Scheduler.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	//logger.Info(" helo_charge request:", string(bt))
	//request.Ext.Ctx
	rsp, err1 := ctl.DoRequest(context.Background(), http.MethodPost, nameServerURL, nil, bytes.NewReader(bt))
	if err1 != nil {
		return nil, err1
	}

	ret := &param.ChargeResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}

	if ret.Status != param.SUCCESS {
		return nil, errors.New("charge fail")
	}

	return ret, nil
}

func SearchRep(request *dict.UploadFinishOrder) (*param.GetOrderRepResponse, error) {
	nameServerURL := fmt.Sprintf("%s://%s/scheduler/v1/searchRep", config.ServerCfg.Request.Protocol, config.ServerCfg.Scheduler.Url)
	bt, err := json.Marshal(request)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("CheckRepProcessor ,searchRep, json marshal request fail ,request: %v, err: %v", *request, err.Error()))
	}

	rsp, err1 := ctl.DoRequest(context.Background(), http.MethodGet, nameServerURL, nil, bytes.NewReader(bt))
	if err1 != nil {
		return nil, err1
	}

	ret := &param.GetOrderRepResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}

	if ret.Status != param.SUCCESS {
		return nil, errors.New("searchRep replication info  fail")
	}

	return ret, nil
}

func GetReplicationStrategy(orderId string) (*param.StrategyInfo, error) {
	nameServerURL := fmt.Sprintf("%s://%s/strategy/v1/getReplicateStrategy", config.ServerCfg.Request.Protocol, config.ServerCfg.Strategy.Url)
	request := param.GetStrategyRequset{
		OrderId: orderId,
	}

	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	rsp, err1 := ctl.DoRequest(context.Background(), http.MethodGet, nameServerURL, nil, bytes.NewReader(bt))
	if err1 != nil {
		return nil, err1
	}

	ret := &param.GetStrategyResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}

	if ret.Status != param.SUCCESS {
		return nil, errors.New("get strategy fail")
	}

	return ret.Strategy, nil
}

//1.重复的.
//2.重新执行删除， 新的，执行删除，可能别人备份成功了，所以可能删除，可能不删除,
//重复的就没上传，所以不用处理,不重复，但是有新订单，也不用删除, 新订单，如果有重复的cid,则不能删除。
func GetOrderDeleteStrategy(orderId string) (*param.StrategyInfo, error) {
	nameServerURL := fmt.Sprintf("%s://%s/strategy/v1/getOrderDeleteStrategy", config.ServerCfg.Request.Protocol, config.ServerCfg.Strategy.Url)
	request := param.GetStrategyRequset{
		OrderId: orderId,
	}

	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	rsp, err1 := ctl.DoRequest(context.Background(), http.MethodGet, nameServerURL, nil, bytes.NewReader(bt))
	if err1 != nil {
		return nil, err1
	}

	ret := &param.GetStrategyResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}

	if ret.Status != param.SUCCESS {
		return nil, errors.New("get getDeleteReplicateStrategy fail")
	}

	return ret.Strategy, nil
}

func GetFidDeleteStrategy(request *param.GetFidDeleteStrategyRequest) (*param.StrategyInfo, error) {
	nameServerURL := fmt.Sprintf("%s://%s/strategy/v1/getFidDeleteStrategy", config.ServerCfg.Request.Protocol, config.ServerCfg.Strategy.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	rsp, err1 := ctl.DoRequest(request.Ext.Ctx, http.MethodGet, nameServerURL, nil, bytes.NewReader(bt))
	if err1 != nil {
		return nil, err1
	}

	ret := &param.GetStrategyResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}

	if ret.Status != param.SUCCESS {
		return nil, errors.New("get getDeleteReplicateStrategy fail")
	}

	return ret.Strategy, nil
}

func DeleteFid(request *param.DeleteOrderFidRequest) (*param.DeleteOrderFidResponse, error) {
	nameServerURL := fmt.Sprintf("%s://%s/scheduler/v1/deleteFid", config.ServerCfg.Request.Protocol, config.ServerCfg.Scheduler.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	rsp, err := ctl.DoRequest(request.Ext.Ctx, http.MethodPost, nameServerURL, nil, bytes.NewReader(bt))
	if err != nil {
		return nil, err
	}

	ret := &param.DeleteOrderFidResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}

	if ret.Status != param.SUCCESS {
		return nil, errors.New("delete order upload task  fail")
	}

	return ret, nil
}

//
func CreateStrategy(request *param.CreateStrategyRequest) (*param.CreateStrategyResponse, error) {
	url := fmt.Sprintf("%s://%s/strategy/v1/createStrategy", config.ServerCfg.Request.Protocol, config.ServerCfg.Strategy.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	rsp, err1 := ctl.DoRequest(context.Background(), http.MethodPost, url, nil, bytes.NewReader(bt))
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

func Log(level int, name, errInfo string, event interface{}) {
	bt, _ := json.Marshal(event)
	switch level {
	case INFO:
		logger.Infof("%v, fail: %v, event: %v", name, errInfo, string(bt))
	case WARN:
		logger.Warnf("%v, fail: %v, event: %v", name, errInfo, string(bt))
	case ERROR:
		logger.Errorf("%v, fail: %v, event: %v", name, errInfo, string(bt))
	}
}
