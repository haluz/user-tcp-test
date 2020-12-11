package client

import (
	"bufio"
	"fmt"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

type Client struct {
	StopChan chan struct{}
}

func (c *Client) Start(userID int, wait bool) {
	logrus.WithField("userID", userID).Info("user client start")

	var wg sync.WaitGroup

	conn, err := net.Dial("tcp", "localhost:7777")
	if err != nil {
		logrus.WithError(err).Error("failed to open connection")
		return
	}

	_, err = conn.Write([]byte(fmt.Sprintf("{\"user_id\": %d}\n", userID)))
	if err != nil {
		logrus.WithError(err).Error("failed to send login request")
		conn.Close()
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				logrus.WithField("userID", userID).WithError(err).Info("error on reading, stop client")
				conn.Close()
				return
			}
			logrus.WithField("userID", userID).Info("message from server: " + message)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-c.StopChan
		conn.Close()
		logrus.WithField("userID", userID).Info("user client stop")
	}()

	if wait {
		wg.Wait()
	}
}

func (c *Client) Stop() {
	c.StopChan <- struct{}{}
}
