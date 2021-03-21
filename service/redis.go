package service

import (
	"github.com/yybirdcf/micro/service/cache"

	"github.com/micro/go-micro/errors"
)

//实例化全局配置的redis
type RedisService struct {
	redisServers map[string]*cache.RedisServer
}

func NewRedisService(srvs map[string]interface{}) *RedisService {
	redisServers := make(map[string]*cache.RedisServer)
	for key, val := range srvs {
		v := val.(map[string]interface{})
		redisServers[key] = cache.NewRedisServer(v)
	}
	return &RedisService{
		redisServers: redisServers,
	}
}
func (srv *RedisService) Close() {
	if srv.redisServers == nil {
		return
	}
	for _, server := range srv.redisServers {
		server.Close()
	}
}
func (srv *RedisService) GetServer(name string) (*cache.RedisServer, error) {
	if server, ok := srv.redisServers[name]; ok {
		return server, nil
	}
	return nil, errors.InternalServerError("service.redis", "redis server not found: %s", name)
}
