package service

import (
	"fmt"

	"github.com/yybirdcf/micro/service/cache"

	"github.com/micro/go-micro/errors"
)

//根据线上生成对应cache实例，可以配置memcache或者redis，或者自定义其它
type MemcacheService struct {
	cacheServers map[string]cache.CacheInter
}

func NewMemcacheService(srvs map[string]interface{}) *MemcacheService {
	cacheServers := make(map[string]cache.CacheInter)

	for key, val := range srvs {
		cacheServers[key] = cache.NewKetamaMemcacheServer(val)
	}

	return &MemcacheService{
		cacheServers: cacheServers,
	}
}

func (srv *MemcacheService) Shutdown() {
	if srv.cacheServers == nil {
		return
	}

	for _, server := range srv.cacheServers {
		server.Close()
	}
}

func (srv *MemcacheService) GetServer(name string) (cache.CacheInter, error) {
	if server, ok := srv.cacheServers[name]; ok {
		return server, nil
	}

	return nil, errors.New(fmt.Sprintf("cache server not found: %s", name))
}
