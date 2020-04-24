package main

import (
	"github.com/gin-gonic/gin"
	"myth/go-essential/app"
	"myth/go-essential/base/rpc/client"
	log "myth/go-essential/log/logc"
	"myth/go-essential/net/rpc/warden"
	"myth/go-example/http_server/common"
	"myth/go-example/http_server/handler"
	"myth/go-example/http_server/manager"
)

func main() {
	p := app.GetMythApp()
	p.Config = &common.Config{}
	p.Run(
		app.WithManager(func(mpp *app.MythApp) app.Manager {
			log.Info("WithManager")
			manager := manager.NewManager()
			return manager
		}),
		app.WithCronTab(func(mpp *app.MythApp) error {
			log.Info("WithCronTab")
			return nil
		}),
		app.WithHttpServer(func(e *gin.Engine, mpp *app.MythApp) error {
			log.Info("WithHttpServer")
			return nil
		}),
		app.WithRpcClient(func(client client.Client, mpp *app.MythApp) error {
			log.Info("WithRpcClient")
			c := client.(*warden.Client)
			log.Init(nil)
			log.Debug("hello start")
			log.Info("hello start")
			log.Error("hello start")
			log.Warn("hello start")
			log.Emerg("hello start")

			//中间件
			i:=0
			c.UseUnary(handler.ExampleAuthFunc(&i))
			hdr := handler.NewHandler(c)

			//for ;;i++{
				hdr.SayHello("hello")
			//}

			return nil
		}),
	)
}
