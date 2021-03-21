package app

import (
	"reflect"

	"github.com/yybirdcf/micro/example/bis"
	"github.com/yybirdcf/micro/example/repository"

	"github.com/yybirdcf/micro/service"

	log "github.com/micro/go-micro/v2/logger"

	"github.com/spf13/viper"
)

var (
	AppCtx *App
)

func init() {
	AppCtx = newApp()
}

//全局app
type App struct {
	//配置
	cfg *viper.Viper

	ServiceRegister    *service.Register
	RepositoryRegister *repository.Register
	BisRegister        *bis.Register
}

func newApp() *App {
	return &App{}
}

//初始化服务
func (app *App) InitApp(cfg *viper.Viper) {
	app.cfg = cfg

	//注册全局服务组件
	app.registerServiceToApp()
	//注册全局repository
	app.registerRepositoryToApp()
	//注册全局bis
	app.registerBisToApp()

}

func (app *App) registerServiceToApp() {
	var err error
	serviceRegister := &service.Register{}
	//初始化基础服务
	app.ServiceRegister.MysqlService, err = service.NewMysqlService(app.cfg.GetStringMap("mysql"))
	if err != nil {
		log.Fatal(err)
	}

	app.ServiceRegister = serviceRegister
}

func (app *App) registerRepositoryToApp() {
	register := &repository.Register{}
	args := []reflect.Value{reflect.ValueOf(app.ServiceRegister)}
	app.reflectCallMethod(register, "InitIns", args)
	app.RepositoryRegister = register
}

func (app *App) registerBisToApp() {
	register := &bis.Register{}
	args := []reflect.Value{reflect.ValueOf(app.ServiceRegister), reflect.ValueOf(register)}
	app.reflectCallMethod(register, "InitIns", args)

	app.BisRegister = register
}

func (app *App) reflectCallMethod(structPointer interface{}, methodName string, args []reflect.Value) {
	bisRefValues := reflect.ValueOf(structPointer)
	if bisRefValues.Kind() != reflect.Ptr || bisRefValues.IsNil() {
		panic("structPointer is not kind of point of struct ")
	}
	for i := 0; i < bisRefValues.Elem().NumField(); i++ {
		field := bisRefValues.Elem().Field(i)
		realFieldAddr := field.Addr().Interface()
		if realFieldAddr == nil {
			continue
		}
		fieldValues := reflect.ValueOf(realFieldAddr)
		method := fieldValues.MethodByName(methodName)
		method.Call(args)
	}
}
