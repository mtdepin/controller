package utils

import (
	"bytes"
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
	WARN = 1
)

func Replicate(request *param.ReplicationRequest) (*param.ReplicationResponse, error) {
	nameServerURL := fmt.Sprintf("%s://%s/scheduler/v1/replicate", config.ServerCfg.Request.Protocol, config.ServerCfg.Scheduler.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	rsp, err1 := ctl.DoRequest(http.MethodPost, nameServerURL, nil, bytes.NewReader(bt))
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
	nameServerURL := fmt.Sprintf("%s://%s/scheduler/v1/delete", config.ServerCfg.Request.Protocol, config.ServerCfg.Scheduler.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	rsp, err1 := ctl.DoRequest(http.MethodPost, nameServerURL, nil, bytes.NewReader(bt))
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
	rsp, err1 := ctl.DoRequest(http.MethodPost, nameServerURL, nil, bytes.NewReader(bt))
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

	rsp, err1 := ctl.DoRequest(http.MethodGet, nameServerURL, nil, bytes.NewReader(bt))
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

	rsp, err1 := ctl.DoRequest(http.MethodGet, nameServerURL, nil, bytes.NewReader(bt))
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

func Log(level int, name, errInfo string, event interface{}) {
	bt, _ := json.Marshal(event)
	logger.Warnf("%v, fail: %v, event: %v", name, errInfo, string(bt))
}
