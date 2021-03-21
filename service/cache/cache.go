package cache

type CacheInter interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, expiration int32) error
	Del(key string) error
	Decr(key string, delta uint64) (uint64, error)
	Incr(key string, delta uint64) (uint64, error)
	Close()
}
