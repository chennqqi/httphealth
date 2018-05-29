package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/chennqqi/goutils/netperf"

	"github.com/NebulousLabs/fastrand"
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
		http.HandleFunc("/"+consul.CONSUL_HEALTH_PATH,
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				hostname, _ := os.Hostname()
				fmt.Fprintf(w, "%v", hostname)
				v := rand.Int() % 1024
				x := fastrand.Bytes(v)
				fmt.Fprintf(w, "%v: %v\n", hostname, hex.EncodeToString(x))
			})
		netperf.ListenAndServe(cfg.HealthHost)
	}()

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
				req := fmt.Sprintf("http://%v:%d/%v", cat.Address, cat.ServicePort, consul.CONSUL_HEALTH_PATH)
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
