package main

import (
	"log"

	"github.com/andreymgn/RSOI-comment/pkg/comment"
	"github.com/andreymgn/RSOI/pkg/tracer"
)

const (
	CommentAppID     = "CommentAPI"
	CommentAppSecret = "PT6RUHLokksaBdIj"
)

func runComment(port int, connString, jaegerAddr, redisAddr, redisPassword string, redisDB int) error {
	tracer, err := tracer.NewTracer("comment", jaegerAddr)
	if err != nil {
		log.Fatal(err)
	}

	knownKeys := map[string]string{CommentAppID: CommentAppSecret}

	server, err := comment.NewServer(connString, redisAddr, redisPassword, redisDB, knownKeys)
	if err != nil {
		log.Fatal(err)
	}
	return server.Start(port, tracer)
}
