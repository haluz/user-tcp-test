package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/haluz/user-notify-test/repo"
	"github.com/haluz/user-notify-test/server"
	"github.com/sirupsen/logrus"
)

func main() {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "INFO"
	}

	l, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Fatalln(err.Error())
	}
	logrus.SetLevel(l)

	repository, err := repo.NewFriendshipRepo()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatalln("failed to get repository")
	}

	ctx, cancel := context.WithCancel(context.Background())
	go shutdown(cancel)

	s := server.NewServer(repository)
	err = s.Run(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatalln("failed to start server")
	}
}

func shutdown(cancelFunc context.CancelFunc) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)

	for {
		sig := <-signalChan
		switch sig {
		case os.Interrupt:
			logrus.Info("SIGINT received, exit")
			cancelFunc()
			return
		}
	}
}
