package api

import (
	"context"
)

type UploadPieceFidRequest struct {
	RequestId string      `json:"request_id"`
	OrderId   string      `json:"order_id"`
	Group     string      `json:"group"`
	Pieces    []*PieceFid `json:"pieces"`
	Ext       *Extend     `json:"ext,omitempty"`
}

type PieceFid struct {
	Fid        string `bson:"fid,omitempty" json:"fid"`
	Rep        int    `bson:"rep,omitempty" json:"rep,omitempty"`
	RepMin     int    `bson:"min_rep,omitempty" json:"min_rep,omitempty"`
	RepMax     int    `bson:"max_rep,omitempty" json:"max_rep,omitempty"`
	Encryption int    `bson:"encryption,omitempty" json:"encryption,omitempty"`
	Expire     uint64 `bson:"expire,omitempty" json:"expire"`
	Level      int    `bson:"level,omitempty" json:"level,omitempty"`
	Size       int    `bson:"size,omitempty" json:"size"`
	Name       string `bson:"name,omitempty" json:"name"`
	Ajust      int    `bson:"ajust,omitempty" json:"ajust,omitempty"`
}

type UploadTaskRequest struct {
	RequestId  string  `json:"request_id"`
	UserId     string  `json:"user_id"`
	UploadType int     `json:"upload_type"`
	PieceNum   int     `json:"piece_num"`
	NodeNum    int     `json:"node_num"`
	Size       uint64  `json:"size"`
	Name       string  `json:"name"`
	RemoteIp   string  `json:"remote_ip"`
	Group      string  `json:"group"`
	NasList    []*Node `json:"nas_list,omitempty"`
	Ext        *Extend `json:"ext,omitempty"`
}

type CreateTaskRequest struct {
	RequestId string  `json:"request_id"`
	Type      int     `json:"task_type"`
	Ext       *Extend `json:"ext,omitempty"`
}

type GetKNodesRequest struct {
	Group   string  `json:"group"`
	NodeNum int     `json:"node_num"`
	Ext     *Extend `json:"ext,omitempty"`
}

type CheckBalanceRequest struct {
	UserId string  `json:"user_id"`
	Ext    *Extend `json:"ext,omitempty"`
}

type Extend struct {
	Ctx context.Context `json:"ctx,omitempty"`
}
