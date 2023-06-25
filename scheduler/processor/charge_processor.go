package processor

import (
	"bytes"
	ctl "controller/pkg/http"
	"controller/scheduler/config"
	e "controller/scheduler/event"
	"controller/scheduler/param"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ChargeProcessor struct {
	pipeLine chan *e.Event
}

func (p *ChargeProcessor) Init(size, num int32) {
	p.pipeLine = make(chan *e.Event, size)

	for i := int32(0); i < num; i++ {
		go p.Handle()
	}
}

func (p *ChargeProcessor) AddEvent(event *e.Event) {
	p.pipeLine <- event
}

func (p *ChargeProcessor) Handle() {
	for {
		msg := <-p.pipeLine
		p.Charge(msg)
	}
}

func (p *ChargeProcessor) Charge(msg *e.Event) {
	request := msg.Data.(*param.ChargeRequest)
	if rsp, err := p.charge(request); err == nil {
		msg.Ret <- rsp
	} else {
		log(WARN, "Charge ", err.Error(), request)
		msg.Ret <- &param.ChargeResponse{Status: param.FAIL}
	}
}

func (p *ChargeProcessor) charge(request *param.ChargeRequest) (*param.ChargeResponse, error) {
	nameServerURL := fmt.Sprintf("%s://%s/account/v1/charge", config.ServerCfg.Request.Protocol, config.ServerCfg.Account.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	rsp, err1 := ctl.DoRequest(http.MethodPost, nameServerURL, nil, bytes.NewReader(bt))
	if err1 != nil {
		return nil, err1
	}

	ret := &param.ChargeResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}

	if ret.Status != param.SUCCESS {
		return nil, errors.New("replicate fail")
	}

	return ret, nil
}
