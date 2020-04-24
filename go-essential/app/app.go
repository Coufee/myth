package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"myth/go-essential/base/rpc/client"
	"myth/go-essential/base/rpc/server"
	"myth/go-essential/conf"
	log "myth/go-essential/log/logc"
	"myth/go-essential/net/rpc/warden"
	"myth/go-essential/utils"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	WorkFlowTypeSync  = 1 //同步执行
	WorkFlowTypeAsync = 2 //异步执行

	WorkFlowNameBase       = "base"
	WorkFlowNameLog        = "log"
	WorkFlowNameHttpServer = "http"
	WorkFlowNameWebSock    = "web_sock"
	WorkFlowNameRpcServer  = "rpc"
	WorkFlowNameManager    = "manager"
	WorkFlowNameCron       = "cron"
)

var (
	_LoadType = conf.LoadConfigTypeIni
	_FilePath = conf.LocalConfigFilePath
	_EtcdAddr = conf.EtcdConfigAddress
)

func init() {
	addFlag(flag.CommandLine)
}

func addFlag(args *flag.FlagSet) {
	args.StringVar(&_LoadType, "conf.load_type", _LoadType, "配置类型")
	args.StringVar(&_FilePath, "conf.file_path", _FilePath, "配置路径")
	args.StringVar(&_EtcdAddr, "conf.etcd_addr", _EtcdAddr, "etcd地址")
}

type Manager interface {
	Start() error
	Close() error
}

type WorkFlow struct {
	Type    int
	Name    string
	Process func(mythApp *MythApp) error
	Close   func(mythApp *MythApp) error
}

type MythApp struct {
	Name          string
	Usage         string
	Version       string
	Config        interface{}
	ConfigWatcher conf.ConfigWatcher
	Manager       Manager
	Upgrade       websocket.Upgrader
	WorkFlows     []WorkFlow
}

var (
	appOnce     sync.Once
	appInstance *MythApp
)

func GetMythApp() *MythApp {
	appOnce.Do(func() {
		appInstance = &MythApp{}
	})

	return appInstance
}

func (mpp *MythApp) CliRun(workflow ...WorkFlow) error {
	_, _ = time.LoadLocation("Asia/Shanghai")
	mpp.WorkFlows = workflow

	app := &cli.App{
		Name:    mpp.Name,
		Usage:   mpp.Usage,
		Version: mpp.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config,c",
				Value: "config",
				Usage: "json config file path",
			},
		},
		Action: func(ctx *cli.Context) error {
			wg := sync.WaitGroup{}
			wg.Add(1)

			for _, wf := range mpp.WorkFlows {
				if wf.Type == WorkFlowTypeSync {
					if err := wf.Process(mpp); err != nil {
						panic(err)
					}
				} else {
					localWorkflow := wf
					wg.Add(1)
					go func() {
						if err := localWorkflow.Process(mpp); err != nil {
							panic(err)
						}
						wg.Done()
					}()
				}
			}

			wg.Done()
			wg.Wait()
			log.Info("Run Myth App All Done")
			return nil
		},
	}

	return app.Run(os.Args)
}

func (mpp *MythApp) Run(workflow ...WorkFlow) error {
	log.Info("Run Myth App All Start")
	_, _ = time.LoadLocation("Asia/Shanghai")
	mpp.WorkFlows = workflow

	//configLoader := &conf.ConfigLoader{
	//	Name:          mpp.Name,
	//	Config:        mpp.Config,
	//	LoadConfigType: _LoadType,
	//	FilePath: _FilePath,
	//	EtcdEndpoint: _EtcdAddr,
	//	ConfigWatcher: mpp.defaultConfigWatcher,
	//}
	//if err := configLoader.Load(); err != nil {
	//	panic(err)
	//	return err
	//}

	wg := sync.WaitGroup{}
	for _, wf := range mpp.WorkFlows {
		if wf.Type == WorkFlowTypeSync {
			if err := wf.Process(mpp); err != nil {
				panic(err)
			}
		} else {
			localWorkflow := wf
			wg.Add(1)
			go func() {
				if err := localWorkflow.Process(mpp); err != nil {
					panic(err)
				}
				wg.Done()
			}()
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("Get A Signal %v", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			mpp.Close()
			log.Info("Myth App Exit")
			time.Sleep(time.Second)
			return nil
		case syscall.SIGHUP:
		default:
			return nil
		}
	}

	return nil
}

func (mpp *MythApp) defaultConfigWatcher(config interface{}) error {
	if utils.VerifyNil(config) {
		return errors.New("defaultConfigWatcher reload config is null ")
	}

	mpp.Config = config
	if !utils.VerifyNil(mpp.ConfigWatcher) {
		mpp.ConfigWatcher(config)
	}
	return nil
}

func (mpp *MythApp) Close() {
	for _, wf := range mpp.WorkFlows {
		if err := wf.Close(mpp); err != nil {
			log.Error("%v server close error(%v)", wf.Name, err)
		}
	}
}

func With(handler func(mpp *MythApp) error) WorkFlow {
	return WorkFlow{
		Type: WorkFlowTypeAsync,
		Name: WorkFlowNameBase,
		Process: func(mythApp *MythApp) error {
			return handler(mythApp)
		},
		Close: func(mythApp *MythApp) error {
			return nil
		},
	}
}

func WithConfig(handler func(mpp *MythApp) error) WorkFlow {
	return WorkFlow{
		Type: WorkFlowTypeAsync,
		Name: WorkFlowNameBase,
		Process: func(mythApp *MythApp) error {
			return handler(mythApp)
		},
		Close: func(mythApp *MythApp) error {
			return nil
		},
	}
}

//func WithLogger() WorkFlow {
//	setLogLevel := func(level string) error {
//		l, err := log.ParseLevel(level)
//		if err != nil {
//			log.Errorf("logger level(%v) error(%v)", viper.GetString("LogLevel"), err)
//			return err
//		}
//
//		log.SetReportCaller(true)
//		log.SetLevel(l)
//		log.SetOutput(os.Stdout)
//		return nil
//	}
//
//	return WorkFlow{
//		Type: WorkFlowTypeSync,
//		Name: WorkFlowNameLog,
//		Process: func(mythApp *MythApp) error {
//			LoggerConfig, ok := mythApp.Config.(conf.GetLogConfig)
//			if !ok {
//				log.Error("config is not TCPServerConfig")
//				return errors.Errorf("config is not TCPServerConfig")
//			}
//
//			mythApp.ConfigWatcher = func(_ interface{}) error {
//				return setLogLevel(LoggerConfig.GetLogConfig().LogLevel)
//			}
//
//			return setLogLevel(LoggerConfig.GetLogConfig().LogLevel)
//		},
//		Close: func(mythApp *MythApp) error {
//			return nil
//		},
//	}
//}

func WithCronTab(handler func(mpp *MythApp) error) WorkFlow {
	return WorkFlow{
		Type: WorkFlowTypeSync,
		Name: WorkFlowNameCron,
		Process: func(mythApp *MythApp) error {
			return handler(mythApp)
		},
		Close: func(mythApp *MythApp) error {
			return nil
		},
	}
}

func WithManager(handler func(mpp *MythApp) Manager) WorkFlow {
	return WorkFlow{
		Type: WorkFlowTypeSync,
		Name: WorkFlowNameManager,
		Process: func(mythApp *MythApp) error {
			manager := handler(mythApp)
			mythApp.Manager = manager
			if err := manager.Start(); err != nil {
				return err
			}

			return nil
		},
		Close: func(mythApp *MythApp) error {
			return mythApp.Manager.Close()
		},
	}
}

func WithHttpServer(handler func(e *gin.Engine, mpp *MythApp) error) WorkFlow {
	e := gin.Default()
	return WorkFlow{
		Type: WorkFlowTypeAsync,
		Name: WorkFlowNameHttpServer,
		Process: func(mythApp *MythApp) error {
			config, ok := mythApp.Config.(conf.GetServerConfig)
			if !ok {
				log.Error("WithRpcServer config is not ServerConfig \n %s", mythApp.Config)
				return errors.Errorf(" WithHttpServer config is not ServerConfig")
			}

			address := fmt.Sprintf("%s:%d", config.GetServerConfig().Host, config.GetServerConfig().RpcPort)
			if err := handler(e, mythApp); err != nil {
				return err
			}

			return e.Run(address)
		},
		Close: func(mythApp *MythApp) error {
			return nil
		},
	}
}

func WithRpcClient(handler func(client client.Client, mpp *MythApp) error) WorkFlow {
	client := warden.NewClient()
	return WorkFlow{
		Type: WorkFlowTypeAsync,
		Name: WorkFlowNameRpcServer,
		Process: func(mythApp *MythApp) error {
			if err := handler(client, mythApp); err != nil {
				return err
			}

			if err := client.Init(); err != nil {
				return err
			}

			return nil
		},
		Close: func(mythApp *MythApp) error {
			return client.Stop()
		},
	}
}

func WithRpcServer(handler func(server server.Server, mpp *MythApp) error) WorkFlow {
	srv := warden.NewServer()
	return WorkFlow{
		Type: WorkFlowTypeAsync,
		Name: WorkFlowNameRpcServer,
		Process: func(mythApp *MythApp) error {
			config, ok := mythApp.Config.(conf.GetServerConfig)
			if !ok {
				log.Error("WithRpcServer config is not ServerConfig \n %s", mythApp.Config.(conf.GetServerConfig))
				return errors.Errorf("WithRpcServer config is not ServerConfig")
			}

			address := fmt.Sprintf("%s:%d", config.GetServerConfig().Host, config.GetServerConfig().RpcPort)
			lis, err := net.Listen("tcp", address)
			if err != nil {
				return err
			}

			if err := srv.Init(); err != nil {
				return err
			}

			if err := handler(srv, mythApp); err != nil {
				return err
			}

			return srv.Serve(lis)
		},
		Close: func(mythApp *MythApp) error {
			return srv.Stop(context.Background())
		},
	}
}

func WithWebSock(handler func(mpp *MythApp) error) WorkFlow {
	return WorkFlow{
		Type: WorkFlowTypeAsync,
		Name: WorkFlowNameWebSock,
		Process: func(mythApp *MythApp) error {
			mythApp.Upgrade = websocket.Upgrader{
				ReadBufferSize:   1024,
				WriteBufferSize:  1024,
				HandshakeTimeout: 5 * time.Second,
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}
			return handler(mythApp)
		},
		Close: func(mythApp *MythApp) error {
			return nil
		},
	}
}
