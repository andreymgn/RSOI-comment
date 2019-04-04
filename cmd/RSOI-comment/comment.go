package main

import (
	"log"

	"github.com/andreymgn/RSOI-comment/pkg/comment"
)

func runComment(port int, connString string) error {
	server, err := comment.NewServer(connString)
	if err != nil {
		log.Fatal(err)
	}

	return server.Start(port)
}
