package web

import (
	"fmt"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"github.com/phayes/freeport"
	"github.com/zxfishhack/cshell/pkg/controllers"
	"log"
)

func InitWebService() (port int, err error) {
	app := iris.New()

	app.Use(recover.New())
	app.Logger().SetLevel(golog.Levels[golog.InfoLevel].Name)
	loggerCfg := logger.DefaultConfig()
	app.Use(logger.New(loggerCfg))

	port, err = freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}

	mvc.Configure(app.Party("/cshell"), func(m *mvc.Application) {
		m.Handle(&controllers.CShellController{})
	})
	mvc.Configure(app.Party("/"), func(m *mvc.Application) {
		m.Handle(&controllers.ConsoleController{})
	})

	go app.Run(iris.Addr(fmt.Sprintf(":%d", port)),
		iris.WithRemoteAddrHeader("X-Real-Ip"),
		iris.WithoutRemoteAddrHeader("X-Forwarded-For"),
		iris.WithoutPathCorrection,
		iris.WithoutBodyConsumptionOnUnmarshal,
	)

	return
}
