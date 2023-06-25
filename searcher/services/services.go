package services

import (
	"controller/searcher/database"
	"controller/searcher/param"
	"controller/searcher/processor"
)

type Service struct {
	searcher *processor.Searcher
}

func (p *Service) Init(db *database.DataBase) {
	p.searcher = new(processor.Searcher)
	p.searcher.Init(db)
}

func (p *Service) Search(request *param.SearchFileRequest) (interface{}, error) {
	return p.searcher.Search(request)
}
