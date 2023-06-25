package processor

import (
	"controller/searcher/database"
	"controller/searcher/dict"
	"controller/searcher/param"
	"errors"
)

type Searcher struct {
	task *database.Task
}

func (p *Searcher) Init(db *database.DataBase) {
	p.task = new(database.Task)
	p.task.Init(db)
}

func (p *Searcher) Search(request *param.SearchFileRequest) (interface{}, error) {
	if request.OrderId != "" {
		return p.searchByOrderId(request.OrderId)
	}

	if request.FileName != "" {
		return p.searchByName(request.FileName)
	}

	return nil, errors.New("Invaid param")
}

func (p *Searcher) searchByOrderId(orderId string) (interface{}, error) {
	tasks, err := p.task.GetTaskInfoByOrderId(orderId)
	if err != nil {
		return nil, err
	}

	rsp := &param.SearchFileResponse{
		Status: param.SUCCESS,
		Datas:  make([]*dict.TaskInfo, 0, len(tasks)),
	}
	for i, _ := range tasks {
		rsp.Datas = append(rsp.Datas, &tasks[i])
	}

	return rsp, nil
}

func (p *Searcher) searchByName(name string) (interface{}, error) {
	tasks, err := p.task.GetTaskInfoByName(name)
	if err != nil {
		return nil, err
	}

	rsp := &param.SearchFileResponse{
		Status: param.SUCCESS,
		Datas:  make([]*dict.TaskInfo, 0, len(tasks)),
	}
	for i, _ := range tasks {
		rsp.Datas = append(rsp.Datas, &tasks[i])
	}

	return rsp, nil
}
