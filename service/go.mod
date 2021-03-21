module github.com/yybirdcf/micro/service

go 1.13

require (
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/gomodule/redigo v1.8.4
	github.com/jinzhu/gorm v1.9.16
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v2 v2.9.1
	github.com/yybirdcf/micro/common/algos v1.0.0
)

replace github.com/yybirdcf/micro/common/algos => ../common/algos