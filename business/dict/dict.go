package dict

const (
	SUCCESS = iota + 1
	FAIL
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
	RequestId  string        `bson:"request_id,omitempty" json:"request_id"`
	UserId     string        `bson:"user_id,omitempty" json:"user_id"`
	UploadType int           `bson:"upload_type,omitempty" json:"upload_type"`
	Tasks      []*UploadTask `bson:"tasks,omitempty" json:"tasks"`
	Group      string        `bson:"group,omitempty" json:"group"`
	NasList    []string      `bson:"nas_list,omitempty" json:"nas_list,omitempty"`
	CreateTime int64         `bson:"create_time,omitempty" json:"create_time,omitempty"`
	UpdateTime int64         `bson:"update_time,omitempty" json:"update_time,omitempty"`
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

type DownloadTask struct {
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
