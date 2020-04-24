package manager

import (
	log "myth/go-essential/log/logc"
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
