package conf

import (
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"myth/go-essential/log/logf"
	"github.com/spf13/viper"
)

func LoadJsonConf(path string,conf interface{}) error {
	if len(path)<=0{
		return errors.New("LoadIniConf path empty")
	}

	//viper.SetEnvPrefix(path)
	viper.AutomaticEnv()
	//replacer := strings.NewReplacer(".", "_")
	//viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigName(path)
	viper.AddConfigPath("./")
	viper.AddConfigPath(".")

	//viper.SetEnvPrefix(path)
	//viper.AutomaticEnv()
	//viper.SetConfigName(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("ConfigLoader ReadConfig failed, err: %s", err)
		return err
	}

	if err := viper.Unmarshal(conf); err != nil {
		log.Errorf("ConfigLoader Unmarshal failed %s", err)
		return err
	}

	configData, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
	log.Debugf("ConfigLoader LoadConfig success, \n%s", configData)

	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Infof("ConfigLoader OnConfigChange op:%s, name:%s", in.Op, in.Name)

		if err := viper.Unmarshal(conf); err != nil {
			log.Errorf("ConfigLoader OnConfigChange failed %s", err)
		} else {
			//if loader.ConfigWatcher != nil {
			//	loader.ConfigWatcher(conf)
			//}
		}
	})

	return nil
}

