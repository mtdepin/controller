package dict

const (
	SUCCESS = iota + 1
	FAIL
)

const (
	LOW = iota
	MIDDLE
	HIGH
)

//region type
const (
	NORMAL_REGION = 0
	CENTER_REGION = 1
)

const (
	LOW_REP_NUM  = 1 //2   //当前测试版本，设置备份数为1.
	MID_REP_NUM  = 2 //3
	HIGH_REP_NUM = 3 //4
)

const (
	TASK_INIT = iota + 1
	TASK_UPLOAD_SUC
	TASK_BEGIN_REP
	TASK_REP_SUC
	TASK_DOWNLOAD_SUC
	TASK_CHARGE_SUC
	TASK_DEL_SUC
	TASK_UPLOAD_FAIL
	TASK_DOWNLOAD_FAIL
	TASK_REP_FAIL
	TASK_DEL_FAIL
	TASK_CHARGE_FAIL
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
	Level      int    `bson:"level,omitempty" json:"level,omitempty"` //0.普通集群, 1.中心集群.
	Status     int    `bson:"status,omitempty" json:"status"`         //1. 有效， 2.删除
	CreateTime int64  `bson:"create_time,omitempty" json:"create_time"`
	UpdateTime int64  `bson:"update_time,omitempty" json:"update_time"`
}

type Rep struct {
	Region     string  `bson:"region,omitempty" json:"region"`
	RealRep    int     `bson:"real_rep,omitempty" json:"real_rep"`
	MinRep     int     `bson:"min_rep,omitempty" json:"min_rep"`
	MaxRep     int     `bson:"max_rep,omitempty" json:"max_rep"`
	RealMinRep int     `bson:"real_min_rep,omitempty" json:"real_min_rep"`
	RealMaxRep int     `bson:"real_max_rep,omitempty" json:"real_max_rep"`
	Expire     uint64  `bson:"expire,omitempty" json:"expire"` //过期时间
	Status     int     `bson:"status,omitempty" json:"status"`
	Weight     float64 `bson:"weight,omitempty" json:"weight"`           //根据开始时间,下载请求并发量等参数设置预估权重,以减少备份次数，实现惰性备份，惰性删除。
	Used       int     `bson:"used,omitempty" json:"used"`               //0:表示此记录未被删除，可以使用, 1: 表示此记录被占用不能删除
	CreateTime int64   `bson:"create_time,omitempty" json:"create_time"` //开始时间, end time = begintime + expire
	UpdateTime int64   `bson:"update_time,omitempty" json:"update_time"`
}

type FidInfo struct { //min,max, 3,4.
	Fid        string                     `bson:"fid,omitempty" json:"fid"`
	Cid        string                     `bson:"cid,omitempty" json:"cid"`         //默认为空，非空有效.
	Origins    string                     `bson:"origins,omitempty" json:"origins"` //文件上传节点
	Region     string                     `bson:"region,omitempty" json:"region"`
	Reps       map[string]map[string]*Rep `bson:"reps,omitempty" json:"reps"` //key:region, orderId, value:备份详情,
	Status     int                        `bson:"status,omitempty" json:"status"`
	CreateTime int64                      `bson:"create_time,omitempty" json:"create_time"`
	UpdateTime int64                      `bson:"update_time,omitempty" json:"update_time"`
}
