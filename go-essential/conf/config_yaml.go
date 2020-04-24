package conf

import (
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	//"github.com/pkg/errors"
	log "myth/go-essential/log/logc"
	"github.com/spf13/viper"
)

func LoadYamlConf(path string,conf interface{}) error {
	if len(path) > 0 {
		viper.SetConfigName(path)
	} else {
		viper.SetConfigName("config")
	}

	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Error("ConfigLoader ReadConfig failed, err: %s", err)
		return err
	}

	if err := viper.Unmarshal(conf); err != nil {
		log.Error("ConfigLoader Unmarshal failed %s", err)
		return err
	}

	configData, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
	log.Debug("ConfigLoader LoadConfig success, \n%s", configData)

	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Info("ConfigLoader OnConfigChange op:%s, name:%s", in.Op, in.Name)

		if err := viper.Unmarshal(conf); err != nil {
			log.Error("ConfigLoader OnConfigChange failed %s", err)
		} else {
			//if loader.ConfigWatcher != nil {
			//	loader.ConfigWatcher(conf)
			//}
		}
	})

	return nil
}

