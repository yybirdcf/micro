package bis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/yybirdcf/micro/example/model"
	"github.com/yybirdcf/micro/example/repository"
	"github.com/yybirdcf/micro/service/cache"

	"github.com/jinzhu/copier"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/yybirdcf/micro/common/algos"
	time2 "github.com/yybirdcf/micro/common/time"
	"github.com/yybirdcf/micro/service"
)

type Base struct {
	serviceRegister    *service.Register
	repositoryRegister *repository.Register
	bisRegister        *Register
}

func (ins *Base) InitIns(serviceRegister *service.Register, repositoryRegister *repository.Register, bisRegister *Register) {
	ins.serviceRegister = serviceRegister
	ins.repositoryRegister = repositoryRegister
	ins.bisRegister = bisRegister
}

func (c *Base) setCache(key string, data interface{}, ttl time.Duration) bool {
	if c.serviceRegister.MemcacheService == nil {
		return true
	}

	srv, err := c.serviceRegister.MemcacheService.GetServer("servers")
	if err != nil {
		log.Errorf("Base setCache failed: %s", err)
		return false
	}

	bs, err := json.Marshal(data)
	if err != nil {
		log.Errorf("Base setCache json err: %s", err)
		return false
	}

	err = srv.Set(key, bs, int32(ttl.Seconds()))
	if err != nil {
		log.Errorf("Base setCache set cache err: %s", err)
		return false
	}

	return true
}

//尝试从缓存里面获取数据，否则从方法里面获取数据并存入缓存，增加了防止缓存穿透，或者缓存击穿，以及缓存雪崩的能力
func (c *Base) getCacheWithDao(key string, ttl time.Duration, randTtl bool, daoFunc func() (model.Model, error), data model.Model, err error) {
	if algos.CacheFilterIns.Exist(key) {
		//防止对数据服务穿透
		err = errors.NotFound("cache.notfound", "cache filter found %s", key)
		return
	}

	srv, err := c.serviceRegister.MemcacheService.GetServer("servers")
	if err != nil {
		log.Errorf("Base getCacheWithDao failed: %s", err)
		return
	}

	//尝试cache
	bs, err := srv.Get(key)
	if err == nil && string(bs) != "null" {
		err = json.Unmarshal(bs, data)
		if err == nil {
			return
		} else {
			log.Errorf("Base getCacheWithDao json.Unmarshal failed: %s", err)
		}
	}

	//请求锁定资源
	ok := algos.DoOnce.Request(key)
	if !ok {
		//等待
		algos.DoOnce.Wait(key)
		//从缓存读取资源，不再对daoFunc发起请求
		bs, err := srv.Get(key)
		if err == nil {
			err = json.Unmarshal(bs, data)
			if err == nil {
				return
			}
		}
		return
	}

	//释放锁定资源
	defer algos.DoOnce.Release(key)

	log.Info("Base getCacheWithDao", key)
	ret, err := daoFunc()
	if err != nil {
		//如果访问后端数据服务出现没有数据错误，增加防止缓存穿透能力，避免对数据服务产生攻击
		if errors.FromError(err).Code == 404 {
			algos.CacheFilterIns.Set(key, "")
		}
		return
	}
	//缓存
	if randTtl {
		//防止缓存雪崩
		c.setCache(key, ret, time2.GetRandomTTL(ttl))
	} else {
		c.setCache(key, ret, ttl)
	}

	copier.Copy(data, ret)
	//data = ret
	return
}

func (c *Base) getCacheListWithDao(ctx context.Context, key string, ttl time.Duration, randTtl bool, daoFunc func() ([]model.Model, error), data interface{}, err error) {
	if algos.CacheFilterIns.Exist(key) {
		//防止对数据服务穿透
		err = errors.NotFound("cache.notfound", "cache filter found %s", key)
		return
	}

	var srv cache.CacheInter
	if c.serviceRegister.MemcacheService != nil {
		srv, err = c.serviceRegister.MemcacheService.GetServer("servers")
		if err != nil {
			log.Errorf("Base getCacheWithDao failed: %s", err)
			return
		}

		//尝试cache
		bs, err := srv.Get(key)
		if err == nil && string(bs) != "null" {
			err = json.Unmarshal(bs, data)
			if err == nil {
				return
			} else {
				log.Errorf("Base getCacheListWithDao json.Unmarshal failed: %s", err)
			}
		}
	}

	//请求锁定资源
	ok := algos.DoOnce.Request(key)
	if !ok {
		//等待
		algos.DoOnce.Wait(key)
		if srv != nil {
			//从缓存读取资源，不再对daoFunc发起请求
			bs, err := srv.Get(key)
			if err == nil {
				err = json.Unmarshal(bs, data)
				if err == nil {
					return
				}
			}
			return
		}
	}

	//释放锁定资源
	defer algos.DoOnce.Release(key)

	log.Info("Base getCacheListWithDao", key)
	ret, err := daoFunc()
	if err != nil {
		//如果访问后端数据服务出现没有数据错误，增加防止缓存穿透能力，避免对数据服务产生攻击
		if errors.FromError(err).Code == 404 {
			algos.CacheFilterIns.Set(key, "")
		}
		return
	}

	//缓存
	if randTtl {
		//防止缓存雪崩
		c.setCache(key, ret, time2.GetRandomTTL(ttl))
	} else {
		c.setCache(key, ret, ttl)
	}

	//这块暂时不知道怎么处理，没法直接赋值，蛋疼
	//data = ret
	bs, err := json.Marshal(ret)
	if err != nil {
		log.Errorf("Base getCacheListWithDao json.Marshal failed: %s", err)
		return
	}
	err = json.Unmarshal(bs, data)

	return
}

func (c *Base) cleanCache(key string) error {
	if c.serviceRegister.MemcacheService == nil {
		return nil
	}
	srv, err := c.serviceRegister.MemcacheService.GetServer("servers")
	if err != nil {
		log.Errorf("Base cleanCache failed: %s", err)
		return nil
	}

	return srv.Del(key)
}
