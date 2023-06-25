package dict

const (
	TASK_INIT = iota + 1
	TASK_UPLOAD_SUC
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

type OrderInfo struct {
	OrderId    string `bson:"order_id,omitempty" json:"order_id"`
	RequestId  string `bson:"request_id,omitempty" json:"request_id"`
	OrderType  int    `bson:"order_type,omitempty" json:"order_type"`
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
	RequestId  string        `bson:"request_id,omitempty" json:"request_id"`
	UserId     string        `bson:"user_id,omitempty" json:"user_id"`
	UploadType int           `bson:"upload_type,omitempty" json:"upload_type"`
	Tasks      []*UploadTask `bson:"tasks,omitempty" json:"tasks"`
	Group      string        `bson:"group,omitempty" json:"group"`
	NasList    []string      `bson:"nas_list,omitempty" json:"nas_list,omitempty"`
	CreateTime int64         `bson:"create_time,omitempty" json:"create_time"`
	UpdateTime int64         `bson:"update_time,omitempty" json:"update_time"`
}

type Task struct {
	Fid     string          `bson:"fid,omitempty" json:"fid"`
	Cid     string          `bson:"cid,omitempty" json:"cid"`
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
	Tasks      map[string]*Task `bson:"task,omitempty" json:"task"` //key:fid, value:备份详情
	Status     int              `bson:"status,omitempty" json:"status"`
	CreateTime int64            `bson:"create_time,omitempty" json:"create_time"`
	UpdateTime int64            `bson:"update_time,omitempty" json:"update_time"`
}

//
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

/*type RepInfo struct{
	num 	   int64 //备份数
	createTime int64
	updateTime int64
}

map[region]*RepInfo

fid  cid  replicate   status  createTime updateTime
*/

type RepInfo struct {
	Region     string `bson:"region,omitempty" json:"region"`
	VirtualRep int    `bson:"virtual_rep,omitempty" json:"virtual_rep"`
	RealRep    int    `bson:"real_rep,omitempty" json:"real_rep"`
	MinRep     int    `bson:"min_rep,omitempty" json:"min_rep"`
	MaxRep     int    `bson:"max_rep,omitempty" json:"max_rep"`
	Status     int    `bson:"status,omitempty" json:"status"`
	CreateTime int64  `bson:"create_time,omitempty" json:"create_time"`
	UpdateTime int64  `bson:"update_time,omitempty" json:"update_time"`
}

type FidInfo struct {
	Fid        string              `bson:"fid,omitempty" json:"fid"`
	Cid        string              `bson:"cid,omitempty" json:"cid"`
	Rep        map[string]*RepInfo `bson:"rep,omitempty" json:"rep"` //key:region, value:备份详情
	Status     int                 `bson:"status,omitempty" json:"status"`
	CreateTime int64               `bson:"create_time,omitempty" json:"create_time"`
	UpdateTime int64               `bson:"update_time,omitempty" json:"update_time"`
}
