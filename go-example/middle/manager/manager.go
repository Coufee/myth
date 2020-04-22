package manager

import (
	log "github.com/sirupsen/logrus"
)

type Manger struct {
}

func NewManager() *Manger {
	return &Manger{}
}

func (manager *Manger) Start() error {
	log.Info("manager Start")

	return nil
}

func (manager *Manger) Close() error {
	log.Info("manager close")
	return nil
}
