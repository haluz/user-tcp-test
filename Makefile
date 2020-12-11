user1:
	WAIT=true USER_ID=1 go run ./cmd/client/client.go

user2:
	WAIT=true USER_ID=2 go run ./cmd/client/client.go

user3:
	WAIT=true USER_ID=3 go run ./cmd/client/client.go

user4:
	WAIT=true USER_ID=4 go run ./cmd/client/client.go

.PHONY: server
server:
	go run ./cmd/server/server.go

test_service:
	go run -race ./test/test.go
