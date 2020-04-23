package conf

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/fsnotify/fsnotify"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"myth/go-essential/log/logf"
	"github.com/spf13/viper"
	"time"
)

const (
	LocalConfigFilePath = "config"
	EtcdConfigAddress   = "127.0.0.1:2379"

	LoadConfigTypeLocal  = "local"
	LoadConfigTypeEtcd   = "etcd"
	LoadConfigTypeZK     = "zookeeper"
	LoadConfigTypeConsul = "consul"
	LoadConfigTypeEureka = "eureka"
)

var LoadConfigMap = map[string]interface{}{
	LoadConfigTypeLocal:  nil,
	LoadConfigTypeEtcd:   nil,
	LoadConfigTypeZK:     nil,
	LoadConfigTypeConsul: nil,
	LoadConfigTypeEureka: nil,
}

type LogConfig struct {
	LogLevel string
}

type GetLogConfig interface {
	GetLogConfig() LogConfig
}

type ServerConfig struct {
	Host                     string
	RpcPort                  int
	HttpPort                 int
	AccessControlAllowOrigin []string
	AccessControlAllowMethod []string
}

type GetServerConfig interface {
	GetServerConfig() ServerConfig
}

type ConfigWatcher func(config interface{}) error

type ConfigLoader struct {
	LoadConfigType string
	Name           string
	FilePath       string
	Config         interface{}
	ConfigWatcher  ConfigWatcher
	EtcdEndpoint   []string
}

func NewConfigLoader(name, filepath string, etcdEndpoint []string) *ConfigLoader {
	return &ConfigLoader{
		Name:         name,
		FilePath:     filepath,
		EtcdEndpoint: etcdEndpoint,
	}
}

func (loader *ConfigLoader) Load() error {
	switch loader.LoadConfigType {
	case LoadConfigTypeLocal:
		return loader.LoadLocalConfig()
	case LoadConfigTypeEtcd:
		return loader.LoadEtcdConfig()
	case LoadConfigTypeZK:
		return loader.LoadZkConfig()
	case LoadConfigTypeConsul:
		return loader.LoadConsulConfig()
	case LoadConfigTypeEureka:
		return loader.LoadEurekaConfig()
	default:
		return errors.New("load config error: config not exist")
	}
	return loader.LoadLocalConfig()
}

//支持文件格式 "json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl"
func (loader *ConfigLoader) LoadLocalConfig() error {
	if len(loader.FilePath) > 0 {
		viper.SetConfigName(loader.FilePath)
	} else {
		viper.SetConfigName("config")
	}

	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("ConfigLoader ReadConfig failed, err: %s", err)
		return err
	}

	if err := viper.Unmarshal(loader.Config); err != nil {
		log.Errorf("ConfigLoader Unmarshal failed %s", err)
		return err
	}

	configData, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
	log.Debugf("ConfigLoader LoadConfig success, \n%s", configData)

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Infof("ConfigLoader OnConfigChange op:%s, name:%s", in.Op, in.Name)

		if err := viper.Unmarshal(loader.Config); err != nil {
			log.Errorf("ConfigLoader OnConfigChange failed %s", err)
		} else {
			if loader.ConfigWatcher != nil {
				loader.ConfigWatcher(loader.Config)
			}
		}
	})

	return nil
}

func (loader *ConfigLoader) LoadZkConfig() error {
	return errors.New("zookeeper config not active")
}

func (loader *ConfigLoader) LoadConsulConfig() error {
	return errors.New("consul config not active")
}

func (loader *ConfigLoader) LoadEurekaConfig() error {
	return errors.New("eureka config not active")
}

func (loader *ConfigLoader) LoadEtcdConfig() error {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   loader.EtcdEndpoint,
		DialTimeout: time.Second * 3,
	})

	if err != nil {
		log.Errorf("ConfigLoader LoadEtcdConfig err:%s", err)
		return err
	}

	key := "/poseidon/config/" + loader.Name
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	response, err := etcdClient.KV.Get(ctx, key)
	if err != nil {
		log.Errorf("ConfigLoader LoadEtcdConfig err:%s", err)
		return err
	}

	if response.Count == 0 {
		log.Errorf("ConfigLoader LoadEtcdConfig response empty")
		return errors.Errorf("ConfigLoader LoadEtcdConfig response empty")

	}

	err = toml.Unmarshal(response.Kvs[0].Value, loader.Config)
	if err != nil {
		log.Errorf("ConfigLoader LoadEtcdConfig err:%s", err)
		return err
	}

	viper.SetConfigType("toml")
	if err := viper.ReadConfig(bytes.NewBuffer(response.Kvs[0].Value)); err != nil {
		log.Errorf("ConfigLoader LoadEtcdConfig err:%s", err)
		return err
	}

	configData, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
	log.Debugf("ConfigLoader ReadRemoteConfig success, \n%s", configData)

	watch := etcdClient.Watcher.Watch(context.Background(), key)

	go func() {
		for resp := range watch {
			log.Infof("ConfigLoader ReadRemoteConfig watcher resp on [%s]", resp)
			value := resp.Events[0].Kv.Value

			err = toml.Unmarshal(value, loader.Config)
			if err != nil {
				log.Errorf("ConfigLoader ReadRemoteConfig watcher err:%s", err)
				continue
			}

			configData, _ := json.MarshalIndent(loader.Config, "", "  ")
			log.Debugf("ConfigLoader ReadRemoteConfig watcher success, \n%s", configData)

			if err := viper.ReadConfig(bytes.NewBuffer(value)); err != nil {
				log.Errorf("ConfigLoader ReadRemoteConfig watcher err:%s", err)
				continue
			}

			if loader.ConfigWatcher != nil {
				loader.ConfigWatcher(loader.Config)
			}
		}
	}()
	return nil
}

//func (loader *ConfigLoader) PushEtcdConfig(path string) error {
//	data, err := ioutil.ReadFile(path)
//	if err != nil {
//		log.Errorf("ConfigLoader PushEtcdConfig read file err:%s, path:%s", err, path)
//		return err
//	}
//
//	etcdClient, err := clientv3.New(clientv3.Config{
//		Endpoints:   loader.EtcdEndpoint,
//		DialTimeout: time.Second * 3,
//	})
//	if err != nil {
//		log.Errorf("ConfigLoader PushEtcdConfig err:%s", err)
//		return err
//	}
//
//	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
//
//	configMap := make(map[string]interface{})
//
//	ext := filepath.Ext(path)
//	if ext == ".yaml" {
//		err = yaml.Unmarshal(data, &configMap)
//	} else {
//		err = json.Unmarshal(data, &configMap)
//	}
//
//	if err != nil {
//		return err
//	}
//
//	configMap[".config_format"] = ext
//
//	for k, v := range configMap {
//		key := fmt.Sprintf("/config/%s/%s", loader.Name, k)
//		value := make([]byte, 0)
//		if ext == ".yaml" {
//			value, err = yaml.Marshal(v)
//			if err != nil {
//				log.Errorf("ConfigLoader LoadEtcdConfig err:%s", err)
//				return err
//			}
//		} else {
//			value, err = json.MarshalIndent(v, "", "  ")
//			if err != nil {
//				log.Errorf("ConfigLoader LoadEtcdConfig err:%s", err)
//				return err
//			}
//		}
//
//		_, err = etcdClient.KV.Put(ctx, key, string(value))
//		if err != nil {
//			log.Errorf("ConfigLoader LoadEtcdConfig err:%s", err)
//			return err
//		}
//	}
//
//	return nil
//}
//
