package utils

import (
	"bytes"
	"context"
	"controller/api"
	"controller/business/config"
	"controller/business/param"
	ctl "controller/pkg/http"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"net/http"
	"time"
)

var c *cache.Cache

const (
	DefaultExpiration = 24 * time.Hour
	CleanupInterval   = 10 * time.Minute
)

func init() {
	c = cache.New(DefaultExpiration, CleanupInterval)
}

func GetUserSecret(userId string, ext *param.Extend) (string, error) {
	//to do 后续添加.
	//return "", nil
	secret, ok := c.Get(userId)
	if ok {
		return secret.(string), nil
	} else {
		rsp, err := getUserInfo(&param.UserInfoRequest{UserId: userId, Ext: ext})
		if err != nil {
			return "", err
		}
		c.Set(userId, rsp.Secret, DefaultExpiration)
		return rsp.Secret, nil
	}
}

func getUserInfo(request *param.UserInfoRequest) (*param.UserInfoResponse, error) {
	url := fmt.Sprintf("%s://%s/user/v1/info", config.ServerCfg.Request.Protocol, config.ServerCfg.User.Url)

	bt, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	rsp, err1 := ctl.DoRequest(request.Ext.Ctx, http.MethodGet, url, nil, bytes.NewReader(bt))
	if err1 != nil {
		return nil, err1
	}

	ret := &param.UserInfoResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}
	if ret.Status != param.SUCCESS {
		return nil, errors.New("get userinfo fail")
	}

	return ret, nil
}

func GetRmSpaceInfo(ctx context.Context, region string) (*api.GetRmSpaceResponse, error) {
	nameServerURL := fmt.Sprintf("%s://%s/api/v1/storage/space?origin=%v", config.ServerCfg.Request.Protocol, config.ServerCfg.Prometheus.Url, region)

	rsp, err := ctl.DoRequest(ctx, http.MethodGet, nameServerURL, nil, nil)
	if err != nil {
		return nil, err
	}

	ret := &api.GetRmSpaceResponse{}
	if err := json.Unmarshal(rsp, ret); err != nil {
		return nil, err
	}

	if ret.Status != param.SUCCESS {
		return nil, errors.New("getRmSpaceInfo fail, rsp: " + string(rsp))
	}

	return ret, nil
}
