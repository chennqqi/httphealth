package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/chennqqi/goutils/consul"
	"github.com/chennqqi/goutils/utils"
)

func main() {
	appName := utils.ApplicationName()

	var cfg Config
	app, err := consul.NewAppWithCfg(&cfg, "")
	if err != nil {
		logrus.Errorf("[%v] NewAppWithCfg error: %v", appName, err)
		return
	}

	listD, _ := time.ParseDuration(cfg.ListInterval)
	reqD, _ := time.ParseDuration(cfg.ReqInterval)
	stopChan := make(chan struct{})
	var stop bool

	go func() {
		m := NewManager()
		m.to = reqD
		defer m.Stop()
		for !stop {
			catlogs, err := app.ListService(appName)
			if err != nil {
				logrus.Errorf("[%v] list service faild %v", appName, err)
				continue
			}
			targets := make([]string, len(catlogs))
			for k, cat := range catlogs {
				req := fmt.Sprintf("http://%v:%d", cat.Address, cat.ServicePort)
				targets[k] = req
			}
			time.Sleep(listD)
			m.Refresh(targets)
		}
	}()

	app.Wait(func(os.Signal) {
		stop = true
		<-stopChan
	}, os.Kill, os.Interrupt)
}
