package dict

const (
	SUCCESS = iota + 1
	FAIL
)

const (
	RM_LFTESPACE_THRESHOLD = 500 * 1024 * 1024 * 1024 //剩余500G
)

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

type DownloadTask struct {
	Fid  string `bson:"fid,omitempty" json:"fid"`
	Cid  string `bson:"cid,omitempty" json:"cid"`
	Name string `bson:"name,omitempty" json:"name"`
	Size int    `bson:"size,omitempty" json:"size"`
}

type DownloadRequestInfo struct {
	RequestId   string          `bson:"request_id,omitempty" json:"request_id"`
	UserId      string          `bson:"user_id,omitempty" json:"user_id"`
	DownlodType int             `bson:"download_type,omitempty" json:"download_type"`
	Tasks       []*DownloadTask `bson:"tasks,omitempty" json:"tasks"`
	CreateTime  int64           `bson:"create_time,omitempty" json:"create_time,omitempty"`
	UpdateTime  int64           `bson:"update_time,omitempty" json:"update_time,omitempty"`
}

type Node struct {
	Address string `bson:"Address,omitempty" json:"Address"`
	Weight  int    `bson:"Weight,omitempty" json:"Weight"`
	RtcId   string `bson:"RtcId,omitempty" json:"RtcId"`
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
