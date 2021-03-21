package bis

import "github.com/yybirdcf/micro/example/model"

type Example struct {
	Base
}

//获取简单信息
func (ins *Example) GetExample(id uint32) (*model.Example, error) {
	return ins.repositoryRegister.Example.Find(id)
}
