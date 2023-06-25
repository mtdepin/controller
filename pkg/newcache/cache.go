package newcache

type Cache interface {
	Get(key string) (val interface{}, exist bool, err error)
	Set(key string, value interface{}, expiration int64) (err error)
	Delete(key string) (err error)
	Scan(cursor uint64, match string, count int64) ([]interface{}, uint64, error)
	Incr(key string, expire int64) (val int64, err error)
	Decr(key string) (val int64, err error)
	SetNX(key string, value interface{}, expiration int64) (flag bool, err error)
}
