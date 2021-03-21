package algos

import "sync"

//防止并发请求同一资源，从而对资源产生压力，比如缓存穿透

type Once struct {
	lock sync.RWMutex
	res  resources
}

type resources map[interface{}]*sync.WaitGroup

var (
	DoOnce = NewOnce()
)

func NewOnce() *Once {
	return &Once{
		lock: sync.RWMutex{},
		res:  make(resources),
	}
}

//请求资源tag，如果没有其它协程对资源tag发起请求，则可以正常发起，并且在请求结束释放资源锁；
//如果已经有协程对资源tag发起请求，则等待，直到其它协程资源请求完成，直接读取返回
func (once *Once) Request(tag interface{}) bool {

	once.lock.Lock()
	defer once.lock.Unlock()

	_, ok := once.res[tag]
	if ok {
		//已经有其他协程请求资源，返回false，进入等待
		return false
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	once.res[tag] = wg

	return true
}

//等待资源tag
func (once *Once) Wait(tag interface{}) {
	once.lock.RLock()
	wg, ok := once.res[tag]
	once.lock.RUnlock()

	if !ok {
		//如果资源已经释放，直接返回，可以获取资源
		return
	}

	//需要等待资源获取
	wg.Wait()
}

//释放资源tag的锁
func (once *Once) Release(tag interface{}) {
	once.lock.Lock()
	defer once.lock.Unlock()

	wg, ok := once.res[tag]
	if !ok {
		return
	}

	wg.Done()
	delete(once.res, tag)
}
