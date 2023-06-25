package distributionlock

import (
	"controller/pkg/newcache"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

//AB之間互斥， AB內部共享
type ABMutexLock struct {
	redis *newcache.RedisClient
}

func NewABMutexLock(addrs []string, password string) *ABMutexLock {
	return &ABMutexLock{
		redis: newcache.NewRedisClient(addrs, password),
	}
}

//region+group
func (p *ABMutexLock) Lock(mutexKey, shareKey string, mutexExpire, shareExpire int64) error {
	pre := time.Now().UnixMilli()
	for true {
		flag, err := p.redis.SetNX(mutexKey, "ok", mutexExpire) //3600s, lock,u
		if err == nil {
			if flag { //设置成功
				//fmt.Printf("setnx success mutexKey: %v, mutexExpire: %v \n", mutexKey, mutexExpire)
				//初始化成1， 防止key存在,  p.redis.Incr(shareKey, shareExpire)
				p.redis.Set(shareKey, 1, shareExpire)

				//fmt.Printf("p.redis.Incr 1 (shareKey: %v , shareExpire: %v, incCount: %v) \n", shareKey, shareExpire, 1)
				return nil
			}

			val, exist, err := p.redis.Get(shareKey)
			if err == nil && exist {
				count, _ := strconv.ParseInt(val.(string), 10, 64)
				if count > 0 {
					p.redis.Incr(shareKey, 0)
					//fmt.Printf("p.redis.Incr 2 (shareKey: %v , shareExpire: %v, incCount: %v) \n", shareKey, 0, v)
					return nil
				}
			}
		}

		if time.Now().UnixMilli()-pre > Timeout { //35s 超時
			errInfo := ""
			if err != nil {
				errInfo = err.Error()
			}
			return errors.New(fmt.Sprintf("get lock timeout fail mutexKey: %s, shareKey: %s , errInfo: %s", mutexKey, shareKey, errInfo))
		}

		time.Sleep(time.Duration(10+rand.Intn(20)) * time.Millisecond)
	}

	return nil
}

func (p *ABMutexLock) UnLock(mutexKey, shareKey string) error {
	val, err := p.redis.Decr(shareKey)

	//fmt.Printf("UnLock 1,  mutexkey: %v,  p.redis.Decr(shareKey: %v, count: %v\n", mutexKey, shareKey, val)
	if err != nil {
		return err
	}

	if val > 0 {
		return nil
	}

	//fmt.Printf("UnLock 2, dec  mutexkey: %v,  p.redis.Decr(shareKey: %v, count: %v \n", mutexKey, shareKey, val)
	p.redis.Delete(shareKey) //共享key 删除，mutex key 删除.
	return p.redis.Delete(mutexKey)
}
