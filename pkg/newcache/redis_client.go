package newcache

import (
	"context"
	"controller/pkg/logger"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	count = 10
)

const (
	KeyNotExist = "redis: nil"
)

type RedisClient struct {
	client *redis.ClusterClient
	//client *redis.Client
}

//集群模式
func (p *RedisClient) InitRedisCluster(addrs []string, password string) bool {
	p.client = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        addrs,
		DialTimeout:  200 * time.Millisecond,
		ReadTimeout:  1000 * time.Millisecond,
		WriteTimeout: 1000 * time.Millisecond,
		//MaxRetries:   5,
		PoolSize:     128,
		MaxRedirects: 3,
		//Username:     "admin", //注意，修改搭配配置中
		//Password:     password,
	})

	if p.client == nil {
		panic(errors.New(fmt.Sprintf("init rediscluster fail, addrs : %s", addrs)))
	}

	return true
}

//单点模式
/*func (p *RedisClient) InitRedisClient(addrs []string, password string) bool {
	p.client = redis.NewClient(&redis.Options{
		Addr:     addrs[0],
		Password: password,
		DB:       0,
	})

	if p.client == nil {
		panic(errors.New(fmt.Sprintf("init rediscluster fail, addrs : %s", addrs)))
	}

	return true
}
*/

func NewRedisClient(addrs []string, password string) *RedisClient {
	redisClient := &RedisClient{}
	redisClient.InitRedisCluster(addrs, password)
	return redisClient
}

func (p *RedisClient) Get(key string) (val interface{}, exist bool, err error) {
	for i := 0; i < count; i++ {
		result := p.client.Get(context.Background(), key)
		result.Result()

		err = result.Err()

		if err != nil && err.Error() == KeyNotExist {
			exist = false
			err = nil
			return
		}

		val = result.Val()
		exist = false

		if err == nil && val != "" {
			exist = true
		}

		if err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	return
}

func (p *RedisClient) Set(key string, value interface{}, expiration int64) (err error) {
	for i := 0; i < count; i++ {
		if err = p.client.Set(context.Background(), key, value, time.Duration(expiration)*time.Millisecond).Err(); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	return
}

func (p *RedisClient) Delete(key string) (err error) {
	for i := 0; i < count; i++ {
		if err = p.client.Del(context.Background(), key).Err(); err == nil {
			return
		}
		if err.Error() == KeyNotExist { //key 不存在
			logger.Info("delete fail, key  not exist", key)
			return nil
		}

		time.Sleep(10 * time.Millisecond)
	}
	return
}

/*func (p *RedisClient) Incr(key string) (int64, error) {
	return p.client.Incr(context.Background(), key).Result()
}*/

func (p *RedisClient) Decr(key string) (val int64, err error) {
	for i := 0; i < count; i++ {
		if val, err = p.decr(key); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	return
}

func (p *RedisClient) Scan(cursor uint64, match string, count int64) ([]interface{}, uint64, error) {
	keys, cur, err := p.client.Scan(context.Background(), cursor, match, count).Result()
	if err != nil {
		return nil, cursor, err
	}

	records := make([]interface{}, 0, len(keys))
	for i, _ := range keys {
		if record, er := p.client.Get(context.Background(), keys[i]).Result(); er == nil {
			records = append(records, record)
		}
	}

	//results, er := p.client.MGet(context.Background(), keys...).Result()
	return records, cur, nil
}

func (p *RedisClient) Incr(key string, expire int64) (val int64, err error) {
	for i := 0; i < count; i++ {
		if val, err = p.incr(key, expire); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	return
}

func (p *RedisClient) SetNX(key string, value interface{}, expiration int64) (flag bool, err error) {
	for i := 0; i < count; i++ {
		if flag, err = p.client.SetNX(context.Background(), key, value, time.Duration(expiration)*time.Millisecond).Result(); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	return
}

func (p *RedisClient) incr(key string, expire int64) (int64, error) {
	if expire <= 0 {
		return p.client.Incr(context.Background(), key).Result()
	}

	ctx := context.Background()
	pipe := p.client.TxPipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Duration(expire)*time.Millisecond)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return incr.Val(), err
}

func (p *RedisClient) decr(key string) (int64, error) {
	return p.client.Decr(context.Background(), key).Result()
}
