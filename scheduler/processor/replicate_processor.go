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
	"strconv"
)

type ReplicateProcessor struct {
	pipeLine  chan *e.Event
	domainMap map[string]*dict.DomainInfo
}

func (p *ReplicateProcessor) Init(size, num int32, domainMap map[string]*dict.DomainInfo) {
	p.pipeLine = make(chan *e.Event, size)
	p.domainMap = domainMap

	for i := int32(0); i < num; i++ {
		go p.Handle()
	}
}

func (p *ReplicateProcessor) AddEvent(event *e.Event) {
	p.pipeLine <- event
}

func (p *ReplicateProcessor) Handle() {
	for {
		p.Replicate(<-p.pipeLine)
	}
}

func (p *ReplicateProcessor) Replicate(msg *e.Event) {
	request := msg.Data.(*param.ReplicationRequest)

	response := &param.ReplicationResponse{
		Status:  param.SUCCESS,
		OrderId: request.OrderId,
		Tasks:   make([]*param.TaskResponse, 0, len(request.Tasks)),
	}

	for _, task := range request.Tasks {
		taskResponse := &param.TaskResponse{
			Fid:          task.Fid,
			Cid:          task.Cid,
			RegionStatus: make(map[string]int, len(task.Reps)),
		}

		for _, rep := range task.Reps {
			replicateRequest := p.generateReplicteRequest(request.OrderId, task.Cid, task.Origins, rep, request.Ext)

			domain, ok := p.domainMap[rep.Region]
			if !ok {
				taskResponse.RegionStatus[rep.Region] = param.FAIL
				log(WARN, "ReplicateProcessor  Replicate ", fmt.Sprintf("region: %v not exist", rep.Region), request)
				continue
			}

			if rsp, err := p.replicate(replicateRequest, domain.Url); err == nil {
				taskResponse.RegionStatus[rep.Region] = rsp.Status
			} else {
				taskResponse.RegionStatus[rep.Region] = param.FAIL
				log(WARN, "ReplicateProcessor  Replicate ", err.Error(), request)
			}
		}

		response.Tasks = append(response.Tasks, taskResponse)
	}

	msg.Ret <- response
}

func (p *ReplicateProcessor) replicate(request *param.TaskReplicateRequest, domainUrl string) (*param.TaskReplicateResponse, error) {
	nameServerURL := fmt.Sprintf("%s://%s/api/v0/pins/%s", config.ServerCfg.Request.Protocol, domainUrl, request.Cid)
	queryParam := make(map[string]string, 10)
	queryParam["meta-orderId"] = request.OrderId
	queryParam["origins"] = request.Origins

	if request.RealRep > 0 {
		queryParam["replication"] = strconv.Itoa(request.RealRep)
	}

	queryParam["replication-min"] = strconv.Itoa(request.MinRep)
	queryParam["replication-max"] = strconv.Itoa(request.MaxRep)
	queryParam["encryption"] = strconv.Itoa(request.Encryption)
	//queryParam["user-allocations"] = "12D3KooWPhjvXUtfCu9Sb1aN56UvW817GZEFKZaJrrg4nuacwtko"

	queryParam["expire-in"] = strconv.FormatUint(request.Expire, 10) + "ms"

	rsp, err1 := ctl.DoRequest(request.Ext.Ctx, http.MethodPost, nameServerURL, queryParam, nil)
	if err1 != nil {
		return nil, err1
	}

	ret := &param.TaskReplicateResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}

	if ret.Status != param.SUCCESS {
		return nil, errors.New("replicate order  task  fail")
	}

	return ret, nil
}

func (p *ReplicateProcessor) generateReplicteRequest(orderId, cid, origins string, rep *param.RepInfo, ext *param.Extend) *param.TaskReplicateRequest {
	return &param.TaskReplicateRequest{
		OrderId:    orderId,
		Cid:        cid,
		Origins:    origins,
		RealRep:    rep.RealRep,
		MinRep:     rep.MinRep,
		MaxRep:     rep.MaxRep,
		Expire:     rep.Expire,
		Encryption: rep.Encryption,
		NasList:    []string{},
		Meta:       make(map[string]string, 1),
		Ext:        ext,
	}
}
