package common

import (
	"gopkg.in/ini.v1"
	"myth/go-essential/conf"
	log "myth/go-essential/log/logc"
)

type Config struct {
	LogConfig    conf.LogConfig
	ServerConfig conf.ServerConfig

	sftpConfig
	uploadDir                string
	processingDir            string
	node                     string
	provinceReportAddress    []string
	corperationReportAddress []string
	provinceReportMap        map[string]int
	groupReportMap           map[string]int
	period                   int
	release                  int
	fileCount                int
}

type sftpConfig struct {
	username string
	password string
	addr     [4]string
	srcPath  string
	dstPath  [2]string
}

type GetConfig interface {
	GetConfig() Config
}

func LoadConfig() (*Config, error) {
	dir := "c.conf"
	cfg := ini.Empty()
	cfg, err := ini.Load(dir)
	if err != nil {
		log.Error("LoadConfig error: [%v]", err)
		return nil, err
	}

	conf := &Config{}
	conf.uploadDir = cfg.Section("ooo").Key("upload_dir").String()
	conf.processingDir = cfg.Section("ooo").Key("processing_dir").String()
	conf.node = cfg.Section("ooo").Key("node").String()
	conf.provinceReportAddress = cfg.Section("ooo").Key("province_report_address").Strings(",")
	conf.corperationReportAddress = cfg.Section("ooo").Key("corperation_report_address").Strings(",")
	conf.provinceReportMap = make(map[string]int, 0)
	conf.groupReportMap = make(map[string]int, 0)
	for _, v := range conf.provinceReportAddress {
		conf.provinceReportMap[v] = 1
	}
	for _, v := range conf.corperationReportAddress {
		conf.groupReportMap[v] = 2
	}

	conf.username = cfg.Section("sftp").Key("username").String()
	conf.password = cfg.Section("sftp").Key("password").String()
	addr := cfg.Section("sftp").Key("addr").Strings(",")
	conf.addr[0] = addr[0]
	conf.addr[1] = addr[1]
	conf.addr[2] = addr[2]
	addr = nil
	dstPath := cfg.Section("sftp").Key("dst_path").Strings(",")
	conf.dstPath[0] = dstPath[0]
	conf.dstPath[1] = dstPath[1]
	dstPath = nil

	conf.period, err = cfg.Section("ooo").Key("period").Int()
	if err != nil {
		return nil, err
	}
	log.Info("period ", conf.period)
	conf.release, err = cfg.Section("ooo").Key("release").Int()
	if err != nil {
		return nil, err
	}

	conf.fileCount, err = cfg.Section("ooo").Key("file_count").Int()
	if err != nil {
		return nil, err
	}

	log.Debug("conf ", conf)
	return conf, nil
}

func (c *Config) GetConfig() Config {
	log.Debug("11111")
	return *c
}

func (c *Config) GetLogConfig() conf.LogConfig {
	return c.LogConfig
}

func (c *Config) GetServerConfig() conf.ServerConfig {
	return c.ServerConfig
}
