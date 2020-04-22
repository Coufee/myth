package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"myth/go-essential/app"
	"myth/go-essential/base/rpc/client"
	"myth/go-essential/net/rpc/warden"
	"myth/go-example/http_server/common"
	"myth/go-example/http_server/handler"
	"myth/go-example/http_server/manager"
)

func main() {
	p := app.GetMythApp()
	p.Config = &common.Config{}
	p.Run(
		app.With(func(mpp *app.MythApp) error {
			log.Info("With")
			return nil
		}),
		app.WithManager(func(mpp *app.MythApp) app.Manager {
			log.Info("WithManager")
			manager := manager.NewManager()
			return manager
		}),
		app.WithCronTab(func(mpp *app.MythApp) error {
			log.Info("WithCronTab")
			return nil
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
