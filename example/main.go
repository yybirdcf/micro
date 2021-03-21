package main

import (
	"github.com/yybirdcf/micro/example/app"
	"github.com/yybirdcf/micro/example/handler"
	example "github.com/yybirdcf/micro/example/proto/example"
	"github.com/yybirdcf/micro/example/subscriber"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/spf13/viper"
)

func main() {

	// read config
	viper.SetConfigFile("./config.json") // 指定配置文件路径
	viper.SetConfigName("config")        // 配置文件名称(无扩展名)
	viper.SetConfigType("json")          // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("/etc/micro/")   // 查找配置文件所在的路径
	viper.AddConfigPath("$HOME/.micro")  // 多次调用以添加多个搜索路径
	viper.AddConfigPath(".")             // 还可以在工作目录中查找配置
	err := viper.ReadInConfig()          // 查找并读取配置文件
	if err != nil {                      // 处理读取配置文件的错误
		log.Fatal(err)
	}

	//初始化框架依赖
	app.AppCtx.InitApp(viper.GetViper())

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.example"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	exampleHandler := &handler.Example{
		ServiceRegister: app.AppCtx.ServiceRegister,
		BisRegister:     app.AppCtx.BisRegister,
	}
	example.RegisterExampleHandler(service.Server(), exampleHandler)

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.service.example", service.Server(), new(subscriber.Example))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
