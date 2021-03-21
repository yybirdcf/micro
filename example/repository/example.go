package repository

import (
	"github.com/yybirdcf/micro/example/model"

	"github.com/micro/go-micro/errors"
)

type Example struct {
	Base
}

func (repo *Example) Find(id uint32) (*model.Example, error) {
	db := repo.serviceRegister.MysqlService.ChooseDb("example", "master")
	if db == nil {
		return nil, errors.InternalServerError("repository.Example.Find", "db not found")
	}

	example := &model.Example{}
	example.Id = uint(id)
	if err := db.First(example).Error; err != nil {
		return nil, err
	}
	return example, nil
}
