package time

import (
	"math/rand"
	"time"
)

func GetCurrentMillseconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetCurrentSeconds() int64 {
	return time.Now().Unix()
}

//根据ttl值 加上一定的随机值
func GetRandomTTL(ttl time.Duration) time.Duration {
	if ttl > (time.Hour * 24) {
		return ttl + time.Second*time.Duration(rand.Int63n(3600))
	} else if ttl > time.Hour {
		return ttl + time.Second*time.Duration(rand.Int63n(300))
	} else if ttl > time.Minute {
		return ttl + time.Second*time.Duration(rand.Int63n(60))
	} else if ttl > time.Second {
		return ttl + time.Second*time.Duration(rand.Int63n(10))
	} else {
		return ttl + time.Millisecond*time.Duration(rand.Int63n(100))
	}
}
