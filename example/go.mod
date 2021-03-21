module github.com/yybirdcf/micro/example/v1

go 1.13

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/golang/protobuf v1.4.2
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/micro/v2 v2.9.3 // indirect
	github.com/spf13/viper v1.7.1
	google.golang.org/grpc v1.26.0 // indirect
	google.golang.org/protobuf v1.23.0
	gorm.io/driver/mysql v1.0.5 // indirect
)
