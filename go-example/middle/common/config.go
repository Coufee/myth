package common

type Config struct {
	LogConfig Log `json:"LogConfig"`
}

type Log struct {
	LogLevel string `json:"LogLevel"`
}

