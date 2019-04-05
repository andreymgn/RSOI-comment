package main

import (
	"github.com/andreymgn/RSOI-comment/pkg/comment"
)

func runComment(port int, connString string) error {
	server, err := comment.NewServer(connString)
	if err != nil {
		return err
	}

	return server.Start(port)
}
