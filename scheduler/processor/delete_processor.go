package processor

import (
	ctl "controller/pkg/http"
	"controller/scheduler/config"
	"controller/scheduler/dict"
	e "controller/scheduler/event"
	"controller/scheduler/param"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type DeleteProcessor struct {
	pipeLine  chan *e.Event
	domainMap map[string]*dict.DomainInfo
}

func (p *DeleteProcessor) Init(size, num int32, domainMap map[string]*dict.DomainInfo) {
	p.pipeLine = make(chan *e.Event, size)
	p.domainMap = domainMap

	for i := int32(0); i < num; i++ {
		go p.Handle()
	}
}

func (p *DeleteProcessor) AddEvent(event *e.Event) {
	p.pipeLine <- event
}

func (p *DeleteProcessor) Handle() {
	for {
		p.Delete(<-p.pipeLine)
	}
}

func (p *DeleteProcessor) Delete(msg *e.Event) {
	request := msg.Data.(*param.DeleteOrderRequest)

	response := &param.DeleteOrderResponse{
		Status:  param.SUCCESS,
		OrderId: request.OrderId,
		Tasks:   make([]*param.UploadTask, 0, len(request.Tasks)),
	}

	for _, task := range request.Tasks {
		deleteRequest := &param.DeleteRequest{
			OrderId: request.OrderId,
			Origins: task.Origins,
			Cid:     task.Cid,
		}

		domain, ok := p.domainMap[task.Region]
		if !ok || domain == nil {
			p.setResponseTask(response, task, param.FAIL)
			log(WARN, "DeleteProcessor, delete task ", fmt.Sprintf("region: %s not exist, or domain is nil ", task.Region), request)
			continue
		}

		if rsp, err := p.delete(deleteRequest, domain.Url); err == nil {
			p.setResponseTask(response, task, rsp.Status)
		} else {
			p.setResponseTask(response, task, param.FAIL) //fid, cid, region.
			log(WARN, "DeleteProcessor, delete task ", err.Error(), request)
		}
	}

	msg.Ret <- response
}

func (p *DeleteProcessor) delete(request *param.DeleteRequest, domainUrl string) (*param.DeleteResponse, error) {
	nameServerURL := fmt.Sprintf("%s://%s/api/v0/pins/%s", config.ServerCfg.Request.Protocol, domainUrl, request.Cid)

	rsp, err1 := ctl.DoRequest(request.Ext.Ctx, http.MethodDelete, nameServerURL, nil, nil)
	if err1 != nil {
		return nil, err1
	}

	ret := &param.DeleteResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}

	if ret.Status != param.SUCCESS {
		return nil, errors.New("delete order upload task  fail")
	}

	return ret, nil
}

func (p *DeleteProcessor) setResponseTask(response *param.DeleteOrderResponse, task *param.UploadTask, status int) {
	response.Tasks = append(response.Tasks, &param.UploadTask{
		Fid:    task.Fid,
		Cid:    task.Cid,
		Region: task.Region,
		Status: status,
	})
}
