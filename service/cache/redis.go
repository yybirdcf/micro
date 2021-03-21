package cache

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	log "github.com/micro/go-micro/v2/logger"
)

var (
	REDIS_DEFAULT_TIMEOUT = 5
)

//暂时先封装一个redis连接池
//不要对外暴露redis的原生操作，方便统一管理，监控
//提供一些常用的操作
type RedisServer struct {
	pool    *redis.Pool
	timeout time.Duration
}

func NewRedisServer(cfg map[string]interface{}) *RedisServer {
	host := cfg["host"].(string)
	port := cfg["port"].(string)
	auth := cfg["auth"].(string)
	timeout := cfg["timeout"].(int)
	db := cfg["db"].(int)

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
			if err != nil {
				log.Errorf("dial redis %s:%d failed: %s\n", host, port, err)
				return nil, err
			}

			if auth != "" {
				if _, err := c.Do("AUTH", auth); err != nil {
					c.Close()
					log.Errorf("auth redis %s:%d failed: %s\n", host, port, err)
					return nil, err
				}
			}

			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		MaxIdle:     20,
		IdleTimeout: 240 * time.Second,
	}

	to := REDIS_DEFAULT_TIMEOUT
	if timeout > 0 {
		to = timeout
	}

	return &RedisServer{
		pool:    pool,
		timeout: time.Duration(to) * time.Second,
	}
}

func (srv *RedisServer) Close() {
	if srv.pool != nil {
		srv.pool.Close()
	}
}

func (srv *RedisServer) Get(key string) ([]byte, error) {
	return redis.Bytes(srv.do("GET", key))
}

func (srv *RedisServer) Set(key string, value []byte, expiration int32) error {
	data, err := redis.String(srv.do("SETEX", key, expiration, value))
	if data == "OK" {
		return nil
	}
	return err
}

func (srv *RedisServer) Del(key string) error {
	_, err := redis.String(srv.do("DEL", key))
	return err
}

func (srv *RedisServer) Decr(key string, delta uint64) (uint64, error) {
	n, err := redis.Uint64(srv.do("DECRBY", key, delta))
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (srv *RedisServer) Incr(key string, delta uint64) (uint64, error) {
	n, err := redis.Uint64(srv.do("INCRBY", key, delta))
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (srv *RedisServer) HGet(key string, field string) string {
	data, err := redis.String(srv.do("HGET", key, field))
	if err != nil {
		return ""
	}

	return data
}

func (srv *RedisServer) HGetAll(key string) map[string]string {
	data, err := redis.StringMap(srv.do("HGETALL", key))
	if err != nil {
		return nil
	}

	return data
}

func (srv *RedisServer) HMGet(key string, fields ...string) map[string]string {
	args := []interface{}{key}
	for _, field := range fields {
		args = append(args, field)
	}
	data, err := redis.Strings(srv.do("HMGET", args...))
	if err != nil {
		return nil
	}

	m := make(map[string]string)
	for index, field := range fields {
		m[field] = data[index]
	}

	return m
}

func (srv *RedisServer) HSet(key string, field string, value string) bool {
	data, err := redis.Int64(srv.do("HSET", key, field, value))
	if err != nil {
		return false
	}

	return data > -1
}

func (srv *RedisServer) HMSet(key string, fieldValues ...string) bool {
	args := []interface{}{key}
	for _, fieldValue := range fieldValues {
		args = append(args, fieldValue)
	}
	data, err := redis.String(srv.do("HMSET", args...))
	if err != nil {
		return false
	}

	return data == "OK"
}

func (srv *RedisServer) RPush(key string, value string) bool {
	data, err := redis.Int64(srv.do("RPUSH", key, value))
	if err != nil {
		return false
	}

	return data > 0
}

func (srv *RedisServer) LPop(key string) string {
	data, err := redis.String(srv.do("LPOP", key))
	if err != nil {
		return ""
	}

	return data
}

func (srv *RedisServer) do(cmd string, args ...interface{}) (interface{}, error) {
	conn := srv.pool.Get()
	defer conn.Close()

	data, err := redis.DoWithTimeout(conn, srv.timeout, cmd, args...)
	if err != nil {
		log.Errorf("redis DoWithTimeout err: %s: %s\n", cmd, err)
	}

	return data, err
}
