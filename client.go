package main

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/levigross/grequests"
)

type Client struct {
	session  *grequests.Session
	to       time.Duration
	target   string
	stopChan chan struct{}
	ctx      context.Context
	cancel   context.CancelFunc
	stop     bool
}

func NewClient(target string, to time.Duration) *Client {
	var client Client
	client.target = target
	client.to = to
	client.stopChan = make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	client.ctx = ctx
	client.cancel = cancel

	return &client
}

func (c *Client) Status() bool {
	return !c.stop
}

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 5
)

func (c *Client) Run() {
	ticker := time.NewTicker(c.to)

	defer ticker.Stop()

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: 3 * time.Second,
	}

FOR_LOOP:
	for {
		select {
		case <-ticker.C:
			req, _ := http.NewRequest("GET", c.target, nil)
			resp, err := client.Do(req)
			if err != nil {
				logrus.Errorf("[%v] Get error: %v", c.target, err)
				break FOR_LOOP
			}

			//to reuse http client
			if resp.Body != nil {
				io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
			}
			if resp.StatusCode != 200 {
				logrus.Errorf("[%v] response: %v", c.target, resp.StatusCode)
				break FOR_LOOP
			}

		case <-c.ctx.Done():
			logrus.Infof("[%v] break by user", c.target)
			break FOR_LOOP
		}
	}

	c.stop = true
	close(c.stopChan)
}

func (c *Client) Stop() {
	c.cancel()
	<-c.stopChan
}
