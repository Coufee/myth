package conf

import (
	"fmt"
	"testing"
)

type Log struct {
	LogLevel string `json:"LogLevel"`
}

type A struct {
	LogConfig Log `json:"LogConfig"`
}

func TestLoadJsonConf(t *testing.T) {
	path := "a"
	conf := &A{}
	err := LoadJsonConf(path, conf)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Print(*conf)
}

type B struct {
	LogConfig string
}

func TestLoadTomlConf(t *testing.T) {
	path := "b"
	conf := &B{}
	err := LoadTomlConf(path, conf)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Print(*conf)
}

type D struct {
	LogConfig string
}

func TestLoadYamlConf(t *testing.T) {
	path := "d"
	conf := &D{}
	err := LoadYamlConf(path, conf)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Print(*conf)
}
