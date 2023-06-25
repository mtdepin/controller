package dict

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

const (
	SEARCH_COUNT = 1500     //1500  //test_code
	REP_COUNT    = 50       //
	CHARGE_COUNT = 50       //
	DEL_COUNT    = 50       //
	Duration     = 36000000 //s , 10h
)

const (
	SUCESS = 1
	FAIL   = 0
)

const (
	FID_INIT = iota
	CID_VALID
)

const (
	REPEATE = 1
)

type OrderInfo struct {
	OrderId    string `bson:"order_id,omitempty" json:"order_id"`
	RequestId  string `bson:"request_id,omitempty" json:"request_id"`
	OrderType  int    `bson:"order_type,omitempty" json:"order_type"`
	PieceNum   int    `bson:"piece_num,omitempty" json:"piece_num,omitempty"`
	Status     int    `bson:"status,omitempty" json:"status"`
	Desc       string `bson:"desc,omitempty" json:"desc"`
	CreateTime int64  `bson:"create_time,omitempty" json:"create_time"`
	UpdateTime int64  `bson:"update_time,omitempty" json:"update_time"`
}

type UploadTask struct {
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

type UploadRequestInfo struct {
	RequestId  string  `bson:"request_id,omitempty" json:"request_id"`
	UserId     string  `bson:"user_id,omitempty" json:"user_id"`
	UploadType int     `bson:"upload_type,omitempty" json:"upload_type"`
	PieceNum   int     `bson:"piece_num,omitempty" json:"piece_num"`
	Size       uint64  `bson:"size,omitempty" json:"size"`
	Name       string  `bson:"name,omitempty" json:"name"`
	RemoteIp   string  `bson:"remote_ip,omitempty" json:"remote_ip"`
	Group      string  `bson:"group,omitempty" json:"group"`
	NasList    []*Node `bson:"nas_list,omitempty" json:"nas_list,omitempty"`
	CreateTime int64   `bson:"create_time,omitempty" json:"create_time"`
	UpdateTime int64   `bson:"update_time,omitempty" json:"update_time"`
}

type Task struct {
	Fid     string          `bson:"fid,omitempty" json:"fid"`
	Cid     string          `bson:"cid,omitempty" json:"cid"`
	Repeate int             `bson:"repeate,omitempty" json:"repeate"` //是否是重复订单，重复订单不用在上传，上传失败，删除的时候也不删除。 0, 不重复， 1： 重复.
	Region  string          `bson:"region,omitempty" json:"region"`   //文件上传区域
	Origins string          `bson:"origins,omitempty" json:"origins"` //文件上传节点
	Reps    map[string]*Rep `bson:"reps,omitempty" json:"reps"`       //key region: val: 备份详情
	Status  int             `bson:"status,omitempty" json:"status"`
}

type Rep struct {
	Region     string `bson:"region,omitempty" json:"region"`
	VirtualRep int    `bson:"virtual_rep,omitempty" json:"virtual_rep"`
	RealRep    int    `bson:"real_rep,omitempty" json:"real_rep"`
	MinRep     int    `bson:"min_rep,omitempty" json:"min_rep"`
	MaxRep     int    `bson:"max_rep,omitempty" json:"max_rep"`
	Expire     uint64 `bson:"expire,omitempty" json:"expire"`
	Encryption int    `bson:"encryption,omitempty" json:"encryption"`
	Status     int    `bson:"status,omitempty" json:"status"`
}

type OrderStateInfo struct {
	OrderId    string           `bson:"order_id,omitempty" json:"order_id"`
	OrderType  int32            `bson:"order_type,omitempty" json:"order_type"`
	PieceNum   int              `bson:"piece_num,omitempty" json:"piece_num"`
	Tasks      map[string]*Task `bson:"task,omitempty" json:"task"` //key:fid, value:备份详情
	Status     int              `bson:"status,omitempty" json:"status"`
	CreateTime int64            `bson:"create_time,omitempty" json:"create_time"`
	UpdateTime int64            `bson:"update_time,omitempty" json:"update_time"`
}

type RegionRep struct {
	Region      string `json:"region"`
	CurRep      int    `json:"curRep"`      //当前备份数,
	Status      int    `json:"status"`      //是否备份成功， 0:备份失败,重新备份, 1: 备份成功.
	CheckStatus int    `json:"checkStatus"` //校验是否备份成功。 0:校验备份失败， 1: 校验备份成功
}

type TaskRepInfo struct {
	Fid     string                `json:"fid"`
	Cid     string                `json:"cid"`
	Regions map[string]*RegionRep `json:"regions"`
}

type RepTask struct {
	Fid     string   `json:"fid"`
	Cid     string   `json:"cid"`
	Regions []string `json:"regions"`
}

type UploadFinishOrder struct {
	OrderId string     `json:"order_id"`
	Tasks   []*RepTask `json:"tasks"`
}

type DownloadTask struct {
	Cid    string `bson:"cid,omitempty" json:"cid"`
	Name   string `bson:"name,omitempty" json:"name"`
	Size   int    `bson:"size,omitempty" json:"size"`
	Status int    `bson:"status,omitempty" json:"status"`
}

type DownloadRequestInfo struct {
	RequestId   string          `bson:"request_id,omitempty" json:"request_id"`
	UserId      string          `bson:"user_id,omitempty" json:"user_id"`
	DownlodType int             `bson:"download_type,omitempty" json:"download_type"`
	Tasks       []*DownloadTask `bson:"tasks,omitempty" json:"tasks"`
	CreateTime  int64           `bson:"create_time,omitempty" json:"create_time,omitempty"`
	UpdateTime  int64           `bson:"update_time,omitempty" json:"update_time,omitempty"`
}

type RepInfo struct {
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
	Fid        string                         `bson:"fid,omitempty" json:"fid"`
	Cid        string                         `bson:"cid,omitempty" json:"cid"`         //默认为空，非空有效.
	Origins    string                         `bson:"origins,omitempty" json:"origins"` //文件上传节点
	Region     string                         `bson:"region,omitempty" json:"region"`
	Reps       map[string]map[string]*RepInfo `bson:"reps,omitempty" json:"reps"` //key:region, orderId, value:备份详情
	Status     int                            `bson:"status,omitempty" json:"status"`
	CreateTime int64                          `bson:"create_time,omitempty" json:"create_time"`
	UpdateTime int64                          `bson:"update_time,omitempty" json:"update_time"`
}

type Node struct {
	Address string `bson:"Address,omitempty" json:"Address"`
	Weight  int    `bson:"Weight,omitempty" json:"Weight"`
	RtcId   string `bson:"RtcId,omitempty" json:"RtcId"`
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
