package param

import (
	"controller/strategy/dict"
	"golang.org/x/net/context"
)

const (
	SUCCESS = 1
	FAIL    = 0
)

type GetStrategyRequset struct {
	OrderId string  `json:"order_id"`
	Ext     *Extend `json:"ext,omitempty"`
}

type CreateStrategyRequest struct {
	OrderId string  `json:"order_id"`
	Region  string  `json:"region"`
	Ext     *Extend `json:"ext,omitempty"`
}

type CreateStrategyResponse struct {
	Status int `json:"status"`
}

type StrategyInfo struct {
	Tasks []*dict.Task `json:"tasks"`
}

type GetStrategyResponse struct {
	Status   int           `json:"status"`
	OrderId  string        `json:"order_id"`
	Strategy *StrategyInfo `json:"strategy"`
}

type GetFidDeleteStrategyRequest struct {
	OrderId string          `json:"order_id"`
	Fids    map[string]bool `json:"fids"`
	Ext     *Extend         `json:"ext,omitempty"`
}

type Extend struct {
	Ctx context.Context `json:"ctx,omitempty"`
}

type GetScheduleStrategyRequest struct {
	RequestId string   `json:"request_id"`
	Fids      []string `json:"fids"`
}

type FidInfo struct {
	Fid  string                   `json:"fid"`
	Reps map[string]*dict.RepInfo `json:"reps"` //key:region.
}

type GetScheduleStrategyResponse struct {
	FidInfos []*FidInfo `json:"fid_infos"`
	Status   int        `json:"status"`
}

//订单的 fid 在某些区域备份失败了，切换到其它区域进行备份。
type GetFidRepStrategyRequest struct {
	OrderId  string              `json:"order_id"`
	FidInfos map[string][]string `json:"fid_infos"` //key: fid, value: 备份失败区域.
}

type GetLoadBalanceStrategyRequest struct {
	OrderId string  `json:"order_id"`
	Region  string  `json:"region"` //当前备份区域
	Ext     *Extend `json:"ext,omitempty"`
}
