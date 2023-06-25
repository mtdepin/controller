package dict

const (
	SUCCESS = iota + 1
	FAIL
)

const MINREPTHRESHOLD = 1

const (
	FIDREP_STATUS_INIT = iota + 0
	FIDREP_STATUS_CIDVALID
)

type RepInfo struct {
	Region     string `bson:"region,omitempty" json:"region"`
	VirtualRep int    `bson:"virtual_rep,omitempty" json:"virtual_rep"`
	RealRep    int    `bson:"real_rep,omitempty" json:"real_rep"`
	MinRep     int    `bson:"min_rep,omitempty" json:"min_rep"`
	MaxRep     int    `bson:"max_rep,omitempty" json:"max_rep"`
	Expire     uint64 `bson:"expire,omitempty" json:"expire"`
	Encryption int    `bson:"encryption,omitempty" json:"encryption"`
	Status     int    `bson:"status,omitempty" json:"status"`
}

type Task struct {
	Fid     string              `bson:"fid,omitempty" json:"fid"`
	Cid     string              `bson:"cid,omitempty" json:"cid"`
	Region  string              `bson:"region,omitempty" json:"region"`   //文件上传区域
	Origins string              `bson:"origins,omitempty" json:"origins"` //文件上传节点
	Reps    map[string]*RepInfo `bson:"reps,omitempty" json:"reps"`       //key region: val: 备份详情
	Status  int                 `bson:"status,omitempty" json:"status"`
}

type StrategyInfo struct {
	OrderId    string  `bson:"order_id,omitempty" json:"order_id"`
	Tasks      []*Task `bson:"tasks,omitempty" json:"tasks"`
	Desc       string  `bson:"desc,omitempty" json:"desc"`
	CreateTime int64   `bson:"create_time,omitempty" json:"create_time"`
	UpdateTime int64   `bson:"update_time,omitempty" json:"update_time"`
}

type TaskInfo struct {
	Fid        string `bson:"fid,omitempty" json:"fid"`
	Cid        string `bson:"cid,omitempty" json:"cid"`
	OrderId    string `bson:"order_id,omitempty" json:"order_id"`
	RequestId  string `bson:"request_id,omitempty" json:"request_id"`
	Rep        int    `bson:"rep,omitempty" json:"rep,omitempty"`
	VirtualRep int    `bson:"virtual_rep,omitempty" json:"virtual_rep,omitempty"`
	RepMin     int    `bson:"min_rep,omitempty" json:"min_rep,omitempty"`
	RepMax     int    `bson:"max_rep,omitempty" json:"max_rep,omitempty"`
	Encryption int    `bson:"encryption,omitempty" json:"encryption,omitempty"`
	Expire     uint64 `bson:"expire,omitempty" json:"expire"`
	Level      int    `bson:"level,omitempty" json:"level,omitempty"`
	Size       int    `bson:"size,omitempty" json:"size"`
	Name       string `bson:"name,omitempty" json:"name"`
	Ajust      int    `bson:"ajust,omitempty" json:"ajust,omitempty"`
	Desc       string `bson:"desc,omitempty" json:"desc,omitempty"`
	Status     int    `bson:"status,omitempty" json:"status,omitempty"`
	CreateTime int64  `bson:"create_time,omitempty" json:"create_time,omitempty"`
	UpdateTime int64  `bson:"update_time,omitempty" json:"update_time,omitempty"`
}

type DomainInfo struct {
	Id         uint64 `bson:"id,omitempty" json:"id"`
	Region     string `bson:"region,omitempty" json:"region"`
	Url        string `bson:"url,omitempty" json:"url"`
	Status     int    `bson:"status,omitempty" json:"status"` //1. 有效， 2.删除
	CreateTime int64  `bson:"create_time,omitempty" json:"create_time"`
	UpdateTime int64  `bson:"update_time,omitempty" json:"update_time"`
}

type FidRepInfo struct {
	Region       string `bson:"region,omitempty" json:"region"`
	VirtualRep   int    `bson:"virtual_rep,omitempty" json:"virtual_rep"`
	RealRep      int    `bson:"real_rep,omitempty" json:"real_rep"`
	MinRep       int    `bson:"min_rep,omitempty" json:"min_rep"`
	MaxRep       int    `bson:"max_rep,omitempty" json:"max_rep"`
	Status       int    `bson:"status,omitempty" json:"status"`
	MinThreshold int    `bson:"min_threshold,omitempty" json:"min_threshold"` // 最低阈值：删除备份时，最少保留数量
	CreateTime   int64  `bson:"create_time,omitempty" json:"create_time"`
	UpdateTime   int64  `bson:"update_time,omitempty" json:"update_time"`
}

type FidInfo struct {
	Fid        string                 `bson:"fid,omitempty" json:"fid"`
	Cid        string                 `bson:"cid,omitempty" json:"cid"`
	Rep        map[string]*FidRepInfo `bson:"rep,omitempty" json:"rep"`       //key:region, value:备份详情
	Status     int                    `bson:"status,omitempty" json:"status"` // 0: init   1: cid 有效
	CreateTime int64                  `bson:"create_time,omitempty" json:"create_time"`
	UpdateTime int64                  `bson:"update_time,omitempty" json:"update_time"`
}

//type FidInfo struct {
//	Fid        string                         `bson:"fid,omitempty" json:"fid"`
//	Cid        string                         `bson:"cid,omitempty" json:"cid"`
//	Reps       map[string]map[string]*RepInfo `bson:"reps,omitempty" json:"reps"`     //key:orderId, region, value:备份详情,
//	Status     int                            `bson:"status,omitempty" json:"status"` // 0: init   1: cid 有效
//	CreateTime int64                          `bson:"create_time,omitempty" json:"create_time"`
//	UpdateTime int64                          `bson:"update_time,omitempty" json:"update_time"`
//}
