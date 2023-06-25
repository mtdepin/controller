package utils

import (
	"bytes"
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

func GetUserSecret(userId string) (string, error) {
	secret, ok := c.Get(userId)
	if ok {
		return secret.(string), nil
	} else {
		rsp, err := getUserInfo(&param.UserInfoRequest{UserId: userId})
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

	rsp, err1 := ctl.DoRequest(http.MethodGet, url, nil, bytes.NewReader(bt))
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
