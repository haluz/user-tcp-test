package server

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	Server struct {
		connections map[int]*connection
		fr          friendshipRepo
		mutex       *sync.RWMutex
	}

	friendshipRepo interface {
		Friends(id int) ([]int, error)
	}
)

func NewServer(repo friendshipRepo) *Server {
	return &Server{
		fr:          repo,
		connections: make(map[int]*connection),
		mutex:       &sync.RWMutex{},
	}
}

func (s *Server) Run(ctx context.Context) error {
	localAddr, err := net.ResolveTCPAddr("tcp", "localhost:7777")
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", localAddr)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := l.SetDeadline(time.Now().Add(time.Second)); err != nil {
				return err
			}

			conn, err := l.Accept()
			if err != nil {
				if os.IsTimeout(err) {
					continue
				}
				return err
			}

			go s.handleLogin(conn)
		}
	}
}

func (s *Server) handleLogin(conn net.Conn) {
	logrus.Debug("accepted new connection, waiting for login request")

	t, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		logrus.WithError(err).Error("read login request failed")
		conn.Close()
		return
	}

	var loginReq loginRequest
	if err := json.Unmarshal([]byte(t), &loginReq); err != nil {
		logrus.WithError(err).Error("login failed")
		conn.Close()
		return
	}

	s.handleConnection(conn, loginReq)
}

func (s *Server) handleConnection(conn net.Conn, loginReq loginRequest) {
	c := s.addUserConn(loginReq)
	logrus.WithField("userID", c.userID).Info("user logged in")

	go func() {
		w := bufio.NewWriter(conn)

	loop:
		for {
			select {
			case status := <-c.friendsStatus:
				logrus.WithField("userID", c.userID).Debug("flush message " + status)

				_, err := w.WriteString(status + "\n")
				if err != nil {
					logrus.WithError(err).Error("failed to write status")
				}
				w.Flush()
			case <-c.done:
				break loop
			}
		}

		w.Flush()
		logrus.WithField("userID", c.userID).Debug("finished writing")
	}()

	// looking for connection close
	go func() {
		r := bufio.NewReader(conn)
		for {
			_, err := r.ReadByte()
			if err != nil {
				if err == io.EOF {
					logrus.WithField("userID", c.userID).Debug("closing connection")
					conn.Close()
					s.removeUserConn(c.userID)
					break
				}

				logrus.WithError(err).Error("error reading from client")
			}
		}
		logrus.WithField("userID", c.userID).Debug("finished reading")
	}()
}

func (s *Server) notifyStatus(userID int, status string) {
	friends, err := s.fr.Friends(userID)
	if err != nil {
		logrus.WithError(err).Error("failed to get friends")
		return
	}

	conns := make([]*connection, 0)
	s.mutex.RLock()
	for _, friend := range friends {
		if v, ok := s.connections[friend]; ok {
			conns = append(conns, v)
		}
	}
	s.mutex.RUnlock()

	for _, conn := range conns {
		go conn.sendStatus(userID, status)
	}
}

func (s *Server) addUserConn(loginReq loginRequest) *connection {
	c := &connection{
		userID:        loginReq.UserID,
		friendsStatus: make(chan string),
		done:          make(chan struct{}),
	}

	s.mutex.Lock()
	s.connections[loginReq.UserID] = c
	s.mutex.Unlock()

	s.notifyStatus(loginReq.UserID, "online")

	return c
}

func (s *Server) removeUserConn(userID int) {
	s.mutex.Lock()
	if v, ok := s.connections[userID]; ok {
		delete(s.connections, userID)
		close(v.done)
	}
	s.mutex.Unlock()
	logrus.WithField("userID", userID).Info("user logged out")

	go s.notifyStatus(userID, "offline")
}
