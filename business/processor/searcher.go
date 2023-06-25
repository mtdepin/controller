package processor

import (
	"bytes"
	"controller/business/config"
	"controller/business/param"
	ctl "controller/pkg/http"
	"encoding/json"
	"fmt"
	"net/http"
)

type Searcher struct {
}

func (p *Searcher) Search(request *param.SearchFileRequest) ([]byte, error) {
	nameServerURL := fmt.Sprintf("%s://%s/searcher/v1/search", config.ServerCfg.Request.Protocol, config.ServerCfg.Searcher.Url)
	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return ctl.DoRequest(request.Ext.Ctx, http.MethodGet, nameServerURL, nil, bytes.NewReader(bt))
}
