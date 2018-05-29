package main

import (
	"time"
	"io"
	"io/ioutil"
	"context"

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

	session := grequests.NewSession(&grequests.RequestOptions{
		RequestTimeout: 2 * time.Second,
		DialTimeout:    2 * time.Second,
	})
	client.session = session

	return &client
}

func (c *Client) Status() bool {
	return !c.stop
}

func (c *Client) Run() {
	ticker := time.NewTicker(c.to)
	session := c.session

	defer ticker.Stop()

FOR_LOOP:
	for {
		select {
		case <-ticker.C:
			resp, err := session.Get(c.target, nil)
			if err != nil {
				logrus.Errorf("[%v] Get error: %v", c.target, err)
				break FOR_LOOP
			}
			if resp.Body!=nil{
				io.Copy(ioutil.Discard, resp.Body)
			}
			resp.Close()
			if resp.StatusCode != 200 {
				logrus.Errorf("[%v] response: %v", c.target, resp.StatusCode)
				break FOR_LOOP
			}

		case <-c.ctx.Done():
			logrus.Info("[%v] break by user", c.target)
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
