package server

import (
	"fmt"
)

type connection struct {
	userID        int
	friendsStatus chan string
	done          chan struct{}
}

func (c *connection) sendStatus(id int, status string) {
	msg := fmt.Sprintf("%d is %s", id, status)

	select {
	case c.friendsStatus <- msg:
	case <-c.done:
		return
	}
}
