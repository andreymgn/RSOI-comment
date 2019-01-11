package comment

import (
	"fmt"
	"log"
	"net"

	"github.com/andreymgn/RSOI/services/auth"

	pb "github.com/andreymgn/RSOI-comment/pkg/comment/proto"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements comments service
type Server struct {
	db   datastore
	auth auth.Auth
}

// NewServer returns a new server
func NewServer(connString, addr, password string, dbNum int, knownApps map[string]string) (*Server, error) {
	db, err := newDB(connString)
	if err != nil {
		return nil, err
	}

	tokenStorage, err := auth.NewInternalAPITokenStorage(addr, password, dbNum, knownApps)
	if err != nil {
		return nil, err
	}

	return &Server{db, tokenStorage}, nil
}

// Start starts a server
func (s *Server) Start(port int, tracer opentracing.Tracer) error {
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	pb.RegisterCommentServer(server, s)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return server.Serve(lis)
}

func (s *Server) checkToken(token string) (bool, error) {
	exists, err := s.auth.Exists(token)
	if err != nil {
		return false, status.Error(codes.Internal, "auth error")
	}

	return exists, nil
}
