package cache

import (
	"fmt"
	"micro/common/algos"

	"github.com/micro/go-micro/errors"

	log "github.com/micro/go-micro/v2/logger"

	"github.com/bradfitz/gomemcache/memcache"
)

//一致性hash
type RingMemcacheServer struct {
	servers []map[string]interface{}
	mcs     map[string]*memcache.Client
	nodes   *algos.HashRing
}

func NewRingMemcacheServer(mss []map[string]interface{}) *RingMemcacheServer {
	m := &RingMemcacheServer{
		servers: mss,
		mcs:     make(map[string]*memcache.Client),
	}

	m.genRing(200)

	return m
}

func (m *RingMemcacheServer) Close() {

}

func (m *RingMemcacheServer) genRing(num int) {
	nodes := algos.NewHashRing(num)

	nodesMap := make(map[string]int)
	for _, server := range m.servers {
		host := server["host"].(string)
		port := server["port"].(string)
		weight := server["weight"].(int)
		ser := fmt.Sprintf("%s:%d", host, port)
		nodesMap[ser] = weight
		m.mcs[ser] = memcache.New(ser)
	}

	nodes.AddNodes(nodesMap)
	m.nodes = nodes
}

func (m *RingMemcacheServer) node(key string) (*memcache.Client, error) {
	if c, ok := m.mcs[m.nodes.GetNode(key)]; ok {
		return c, nil
	}

	log.Errorf("RingMemcacheServer Get err: memcache node not found %s\n", key)

	return nil, errors.New("memcache node not found")
}

func (m *RingMemcacheServer) Get(key string) ([]byte, error) {
	node, err := m.node(key)
	if err != nil {
		return nil, err
	}

	item, err := node.Get(key)
	if err != nil {
		if err != memcache.ErrCacheMiss {
			log.Errorf("RingMemcacheServer Get err: %s\n", err)
		}
		return nil, err
	}

	return item.Value, nil
}

//过期时间秒数，0表示不过期
func (m *RingMemcacheServer) Set(key string, value []byte, expiration int32) error {
	node, err := m.node(key)
	if err != nil {
		return err
	}

	item := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: expiration,
	}

	err = node.Set(item)
	if err != nil {
		log.Errorf("RingMemcacheServer Set err: %s\n", err)
	}
	return err
}

func (m *RingMemcacheServer) Del(key string) error {
	node, err := m.node(key)
	if err != nil {
		return err
	}

	err = node.Delete(key)
	if err != nil {
		log.Errorf("RingMemcacheServer Del err: %s\n", err)
	}
	return err
}

func (m *RingMemcacheServer) Decr(key string, delta uint64) (uint64, error) {
	node, err := m.node(key)
	if err != nil {
		return 0, err
	}

	val, err := node.Decrement(key, delta)
	if err != nil {
		log.Errorf("RingMemcacheServer Decr err: %s\n", err)
	}
	return val, err
}

func (m *RingMemcacheServer) Incr(key string, delta uint64) (uint64, error) {
	node, err := m.node(key)
	if err != nil {
		return 0, err
	}

	val, err := node.Increment(key, delta)
	if err != nil {
		log.Errorf("RingMemcacheServer Incr err: %s\n", err)
	}
	return val, err
}

type KetamaMemcacheServer struct {
	servers []map[string]interface{}
	mcs     map[string]*memcache.Client
	nodes   *algos.Continuum
}

func NewKetamaMemcacheServer(mss []map[string]interface{}) *KetamaMemcacheServer {
	m := &KetamaMemcacheServer{
		servers: mss,
		mcs:     make(map[string]*memcache.Client),
	}

	buckets := make([]algos.Bucket, 1)
	for _, server := range m.servers {
		host := server["host"].(string)
		port := server["port"].(string)
		weight := server["weight"].(int)
		ser := fmt.Sprintf("%s:%d", host, port)
		bucket := algos.Bucket{
			Label:  ser,
			Weight: weight,
		}

		buckets = append(buckets, bucket)
		m.mcs[ser] = memcache.New(ser)
	}

	nodes, err := algos.NewContinuum(buckets)
	if err != nil {
		log.Errorf("NewKetamaMemcacheServer err : %s", err)
	}

	m.nodes = nodes
	return m
}

func (m *KetamaMemcacheServer) Close() {

}

func (m *KetamaMemcacheServer) node(key string) (*memcache.Client, error) {
	if c, ok := m.mcs[m.nodes.Hash(key)]; ok {
		return c, nil
	}

	log.Errorf("KetamaMemcacheServer Get err: memcache node not found %s\n", key)

	return nil, errors.New("memcache node not found")
}

func (m *KetamaMemcacheServer) Get(key string) ([]byte, error) {
	node, err := m.node(key)
	if err != nil {
		return nil, err
	}

	item, err := node.Get(key)
	if err != nil {
		if err != memcache.ErrCacheMiss {
			log.Errorf("KetamaMemcacheServer Get err: %s\n", err)
		}
		return nil, err
	}

	return item.Value, nil
}

//过期时间秒数，0表示不过期
func (m *KetamaMemcacheServer) Set(key string, value []byte, expiration int32) error {
	node, err := m.node(key)
	if err != nil {
		return err
	}

	item := &memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: expiration,
	}

	err = node.Set(item)
	if err != nil {
		log.Errorf("KetamaMemcacheServer Set err: %s\n", err)
	}
	return err
}

func (m *KetamaMemcacheServer) Del(key string) error {
	node, err := m.node(key)
	if err != nil {
		return err
	}

	err = node.Delete(key)
	if err != nil {
		log.Errorf("KetamaMemcacheServer Del err: %s\n", err)
	}
	return err
}

func (m *KetamaMemcacheServer) Decr(key string, delta uint64) (uint64, error) {
	node, err := m.node(key)
	if err != nil {
		return 0, err
	}

	val, err := node.Decrement(key, delta)
	if err != nil {
		log.Errorf("KetamaMemcacheServer Decr err: %s\n", err)
	}
	return val, err
}

func (m *KetamaMemcacheServer) Incr(key string, delta uint64) (uint64, error) {
	node, err := m.node(key)
	if err != nil {
		return 0, err
	}

	val, err := node.Increment(key, delta)
	if err != nil {
		log.Errorf("KetamaMemcacheServer Incr err: %s\n", err)
	}
	return val, err
}
