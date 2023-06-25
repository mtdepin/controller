package dict

const (
	SUCCESS = iota + 1
	FAIL
)

type DomainInfo struct {
	Id         uint64 `bson:"id,omitempty" json:"id"`
	Region     string `bson:"region,omitempty" json:"region"`
	Url        string `bson:"url,omitempty" json:"url"`
	Level      int    `bson:"level,omitempty" json:"level,omitempty"` //0.普通集群, 1.中心集群.
	Status     int    `bson:"status,omitempty" json:"status"`         //1. 有效， 2.删除
	CreateTime int64  `bson:"create_time,omitempty" json:"create_time"`
	UpdateTime int64  `bson:"update_time,omitempty" json:"update_time"`
}
