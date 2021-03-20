package bis

import (
	"example/repository"
	"example/service"
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
