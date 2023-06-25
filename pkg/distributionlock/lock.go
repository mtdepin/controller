package distributionlock

import (
	"controller/pkg/newcache"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	Timeout = 10000
	Expire  = 10000 //10s
	//ShareExpire = 3300000 //3000s
	//ABExpire    = 3600000 //3600s
)

type MutexLock struct {
	redis *newcache.RedisClient
}

func NewMutexLock(addrs []string, password string) *MutexLock {
	return &MutexLock{
		redis: newcache.NewRedisClient(addrs, password),
	}
}

func (p *MutexLock) Lock(key string) error {
	pre := time.Now().UnixMilli()

	for true {
		flag, err := p.redis.SetNX(key, "ok", Expire)
		if flag && err == nil {
			return nil
		}

		if time.Now().UnixMilli()-pre > Timeout { //35s 超時
			errInfo := ""
			if err != nil {
				errInfo = err.Error()
			}
			return errors.New(fmt.Sprintf("get lock timeout fail : %s , errInfo: %s", key, errInfo))
		}

		time.Sleep(time.Duration(10+rand.Intn(20)) * time.Millisecond)
	}
	return nil
}

func (p *MutexLock) UnLock(key string) error {
	for i := 0; i < 10; i++ {
		if err := p.redis.Delete(key); err == nil {
			return nil
		}
		time.Sleep(time.Duration(10) * time.Millisecond)
	}
	return errors.New("unlock fail")
}
