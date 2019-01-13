package main

import (
	"log"

	"github.com/andreymgn/RSOI-comment/pkg/comment"
)

const (
	CommentAppID     = "CommentAPI"
	CommentAppSecret = "PT6RUHLokksaBdIj"
)

func runComment(port int, connString, redisAddr, redisPassword string, redisDB int) error {
	knownKeys := map[string]string{CommentAppID: CommentAppSecret}

	server, err := comment.NewServer(connString, redisAddr, redisPassword, redisDB, knownKeys)
	if err != nil {
		log.Fatal(err)
	}
	return server.Start(port)
}
