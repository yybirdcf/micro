package repository

import "example/service"

type Base struct {
	serviceRegister *service.Register
}

func (ins *Base) InitIns(serviceRegister *service.Register) {
	ins.serviceRegister = serviceRegister
}
