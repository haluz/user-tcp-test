package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/haluz/user-tcp-test/client"
	"github.com/haluz/user-tcp-test/repo"
	"github.com/haluz/user-tcp-test/server"
	"github.com/sirupsen/logrus"
)

func main() {
	go runServer()
	time.Sleep(100 * time.Millisecond) // let server start

	ctx, cancel := context.WithCancel(context.Background())

	m := make(map[int]chan struct{}, 4)
	for i := 1; i < 5; i++ {
		m[i] = make(chan struct{})
		go clientRun(ctx, i, m[i])
	}

	for i := 0; i < 1000; i++ {
		ch := m[rand.Intn(4)+1] // get random client channel
		ch <- struct{}{}

		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}

	cancel()
	time.Sleep(100 * time.Millisecond)
}

func clientRun(ctx context.Context, id int, op chan struct{}) {
	run := false
	var c *client.Client

	for {
		select {
		case <-op:
			if run {
				c.Stop()
				c = nil
				run = false
			} else {
				c = &client.Client{
					StopChan: make(chan struct{}),
				}
				c.Start(id, false)
				run = true
			}
		case <-ctx.Done():
			if run {
				c.Stop()
			}
			return
		}
	}
}

func runServer() {
	repository, err := repo.NewFriendshipRepo()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatalln("failed to get repository")
	}

	s := server.NewServer(repository)
	err = s.Run(context.Background())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatalln("failed to start server")
	}
}
