package main

import (
	"github.com/andreymgn/RSOI-comment/pkg/comment"
	"github.com/andreymgn/RSOI/pkg/tracer"
)

func runComment(port int, connString, jaegerAddr string) error {
	tracer, closer, err := tracer.NewTracer("comment", jaegerAddr)
	defer closer.Close()
	if err != nil {
		return err
	}

	server, err := comment.NewServer(connString)
	if err != nil {
		return err
	}

	return server.Start(port, tracer)
}
