package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"myth/go-essential/app"
	"myth/go-essential/base/rpc/client"
	"myth/go-essential/net/rpc/warden"
	"myth/go-example/client/handler"
	"myth/go-example/client/manager"
)

func main() {
	p := app.GetMythApp()
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
			hdr := handler.NewHandler(c)
			hdr.SayHello("hello")
			return nil
		}),
	)
}