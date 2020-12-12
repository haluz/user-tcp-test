package server

type connection struct {
	userID        int
	friendsStatus chan userStatus
	done          chan struct{}
}

func (c *connection) sendStatus(status userStatus) {
	select {
	case c.friendsStatus <- status:
	case <-c.done:
		return
	}
}
