package main

import (
	"github.com/haluz/user-notify-test/client"
	"os"
	"strconv"
)

func main() {
	userID := os.Getenv("USER_ID")
	wait := os.Getenv("WAIT")

	c := &client.Client{}

	id, _ := strconv.ParseInt(userID, 10, 64)
	w, _ := strconv.ParseBool(wait)

	c.Start(int(id), w)
}
