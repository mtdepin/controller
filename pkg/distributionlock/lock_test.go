package distributionlock

import (
	"controller/pkg/logger"
	"controller/pkg/newcache"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

type Rep struct {
	//OrderId    string  `bson:"order_id,omitempty" json:"order_id"`
	Region     string  `bson:"region,omitempty" json:"region"`
	RealRep    int     `bson:"real_rep,omitempty" json:"real_rep"`
	MinRep     int     `bson:"min_rep,omitempty" json:"min_rep"`
	MaxRep     int     `bson:"max_rep,omitempty" json:"max_rep"`
	RealMinRep int     `bson:"real_min_rep,omitempty" json:"real_min_rep"`
	RealMaxRep int     `bson:"real_max_rep,omitempty" json:"real_max_rep"`
	Expire     uint64  `bson:"expire,omitempty" json:"expire"` //过期时间
	Status     int     `bson:"status,omitempty" json:"status"`
	Weight     float64 `bson:"weight,omitempty" json:"weight"`           //根据开始时间,下载请求并发量等参数设置预估权重,以减少备份次数，实现惰性备份，惰性删除。
	CreateTime int64   `bson:"create_time,omitempty" json:"create_time"` //开始时间, end time = begintime + expire
	UpdateTime int64   `bson:"update_time,omitempty" json:"update_time"`
}

type FidInfo struct { //min,max, 3,4.
	Fid        string                     `bson:"fid,omitempty" json:"fid"`
	Cid        string                     `bson:"cid,omitempty" json:"cid"`   //默认为空，非空有效.
	Reps       map[string]map[string]*Rep `bson:"reps,omitempty" json:"reps"` //key:region, orderId, value:备份详情,
	Status     int                        `bson:"status,omitempty" json:"status"`
	Used       int                        `bson:"used,omitempty" json:"used"` //0:表示此记录未被删除，可以使用, 1: 表示此记录被占用不能删除
	CreateTime int64                      `bson:"create_time,omitempty" json:"create_time"`
	UpdateTime int64                      `bson:"update_time,omitempty" json:"update_time"`
}

var mutex *MutexLock
var cache newcache.Cache

func initLock() {
	mutex = NewMutexLock(strings.Split("10.80.7.28:7001,10.80.7.30:7001", ","), "kepler123456")
	cache = newcache.NewRedisClient(strings.Split("10.80.7.28:7001,10.80.7.30:7001", ","), "kepler123456")
}

func TestCache(t *testing.T) {
	initLock()

	val, exist, err := cache.Get("db8e6e12f30434c6ec67adbba68c48bb")
	if err != nil {
		return
	}

	if !exist { //正常情况，应该存在.
		return
	}

	fidInfo := &FidInfo{}
	if err := json.Unmarshal([]byte(val.(string)), fidInfo); err != nil {
		logger.Error(fmt.Sprintf("createFidRepStrategy json.Unmarshal fidInfo fail, err: %v, fidInfo :%v ", err.Error(), val.(string)))
		return
	}
}

func TestLock(t *testing.T) {
	initLock()

	fmt.Printf("begin\n")

	t0 := time.Now().UnixMilli()

	for i := 0; i < 10; i++ {
		key := "key123"
		t1 := time.Now().UnixMilli()
		if err := mutex.Lock(key); err == nil {
			fmt.Printf("lock key123 succ \n")
			//time.Sleep(100 * time.Second)
			t2 := time.Now().UnixMilli()
			mutex.UnLock(key)
			t3 := time.Now().UnixMilli()

			fmt.Printf(" lock_cost: %v, unlock_cost: %v, total_cost: %v \n", t2-t1, t3-t2, t3-t1)

		} else {
			fmt.Printf("lock key123  fail: %v \n", err.Error())
		}
	}

	t5 := time.Now().UnixMilli()

	fmt.Printf(" total_cost: %v \n", t5-t0)

	//DistributeLock_1()
}
