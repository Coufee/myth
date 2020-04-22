package common

import "myth/go-essential/conf"

type Config struct {
	LogConfig    conf.LogConfig
	ServerConfig conf.ServerConfig
}

func (c *Config) GetLogConfig() conf.LogConfig {
	return c.LogConfig
}

func (c *Config) GetServerConfig() conf.ServerConfig {
	return c.ServerConfig
}