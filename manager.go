package main

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

type Manager struct {
	clients map[string]*Client
	lock    sync.Mutex
	to      time.Duration
}

func NewManager() *Manager {
	var m Manager
	m.clients = make(map[string]*Client)
	return &m
}

func (m *Manager) Refresh(targes []string) {
	for _, t := range targes {
		c, exist := m.clients[t]
		if !exist {
			logrus.Infof("[%v] pushed new", t)
			nc := NewClient(t, m.to)
			go nc.Run()
			m.clients[t] = nc
		} else if !c.Status() {
			nc := NewClient(t, m.to)
			logrus.Infof("[%v] pushed old", t)
			go nc.Run()
			m.clients[t] = nc
		}
	}
}

func (m *Manager) Stop() {
	for _, c := range m.clients {
		c.Stop()
	}
	m.clients = make(map[string]*Client)
}
