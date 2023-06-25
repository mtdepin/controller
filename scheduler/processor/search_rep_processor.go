package processor

import (
	"context"
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

type SearchRepProcessor struct {
	pipeLine  chan *e.Event
	domainMap map[string]*dict.DomainInfo
}

func (p *SearchRepProcessor) Init(size, num int32, domainMap map[string]*dict.DomainInfo) {
	p.pipeLine = make(chan *e.Event, size)
	p.domainMap = domainMap

	for i := int32(0); i < num; i++ {
		go p.Handle()
	}
}

func (p *SearchRepProcessor) AddEvent(event *e.Event) {
	p.pipeLine <- event
}

func (p *SearchRepProcessor) Handle() {
	for {
		p.Search(<-p.pipeLine)
	}
}

func (p *SearchRepProcessor) Search(msg *e.Event) {
	request := msg.Data.(*param.UploadFinishOrder)

	response := &param.GetOrderRepResponse{
		Status:  param.SUCCESS,
		OrderId: request.OrderId,
		Tasks:   make(map[string]*param.TaskRepInfo, len(request.Tasks)),
	}

	regionCids := make(map[string]*RegionCid, 10)

	p.convertsearchRequest(request, response, regionCids)

	for region, regionCid := range regionCids {
		dominInfo, ok := p.domainMap[region]
		if !ok {
			log(WARN, "SearchRepProcessor, search ", fmt.Sprintf("region: %v not exits", region), request)
			continue
		}

		rsp, err := p.search(regionCid.Cids, dominInfo.Url)
		if err != nil {
			log(WARN, fmt.Sprintf("SearchRepProcessor, search  region: %v", region), err.Error(), request)
			continue
		}
		//test code
		//bt, _ := json.Marshal(rsp)
		//logger.Warnf("SearchRepProcessor, search , region: %v   request: %v,  response: %v", region, *request, string(bt))
		////
		p.setCidRepResponse(region, response, rsp)
	}

	msg.Ret <- response
}

func (p *SearchRepProcessor) search(cids string, domainUrl string) (*param.RegionRepResponse, error) {

	nameServerURL := fmt.Sprintf("%s://%s/api/v0/status/%s", config.ServerCfg.Request.Protocol, domainUrl, cids)

	rsp, err1 := ctl.DoRequest(context.Background(), http.MethodGet, nameServerURL, nil, nil)
	if err1 != nil {
		return nil, err1
	}

	ret := &param.RegionRepResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}

	if ret.Status != param.SUCCESS {
		return nil, errors.New("search order  region rep fail")
	}

	return ret, nil
}

func (p *SearchRepProcessor) convertsearchRequest(request *param.UploadFinishOrder, response *param.GetOrderRepResponse, regionCids map[string]*RegionCid) {

	for _, task := range request.Tasks {
		//init response info
		response.Tasks[task.Fid] = &param.TaskRepInfo{Fid: task.Fid, Cid: task.Cid, Regions: make(map[string]*param.RegionRep, 10)}

		//generate seach request.
		for _, region := range task.Regions {
			if regionCid, ok := regionCids[region]; ok {
				if _, ok := regionCid.CidMap[task.Cid]; !ok {
					regionCid.CidMap[task.Cid] = true
					regionCid.Cids += "," + task.Cid
				}
			} else {
				regionCid = &RegionCid{
					CidMap: make(map[string]bool, 10),
					Cids:   task.Cid,
				}
				regionCid.CidMap[task.Cid] = true
				regionCids[region] = regionCid
			}
		}
	}
}

func (p *SearchRepProcessor) setCidRepResponse(region string, rsp *param.GetOrderRepResponse, regionRep *param.RegionRepResponse) {
	if rsp == nil || regionRep == nil {
		return
	}

	for _, task := range rsp.Tasks {
		if rep, ok := regionRep.CidRep[task.Cid]; ok {
			task.Regions[region] = &param.RegionRep{
				Region: region,
				CurRep: rep.PinCount,
				Status: rep.Status,
			}
		}
	}
}
