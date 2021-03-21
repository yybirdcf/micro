module github.com/yybirdcf/micro/example

go 1.13

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/golang/protobuf v1.4.2
	github.com/jinzhu/now v1.1.1 // indirect
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v2 v2.9.1
	github.com/spf13/viper v1.7.1
	github.com/yybirdcf/micro/service v0.0.0-20210321062050-e1cf68303f76
	google.golang.org/protobuf v1.23.0
)
