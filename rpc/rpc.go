package rpc

import (
	"fmt"
	"net"
	"time"

	"github.com/guneyin/app-sdk/utils"

	"github.com/guneyin/app-sdk/logger"
	"google.golang.org/grpc"
)

const defaultRequestTimeout = time.Second * 5

type Server struct {
	server *grpc.Server
	port   string
}

func New(port string, timeout time.Duration) *Server {
	if timeout == 0 {
		timeout = defaultRequestTimeout
	}

	return &Server{
		server: grpc.NewServer(
			grpc.ConnectionTimeout(timeout),
		),
		port: port,
	}
}

func (s *Server) Addr() string {
	return fmt.Sprintf("%s:%s", utils.GetOutboundIP().String(), s.port)
}

func (s *Server) Server() *grpc.Server {
	return s.server
}

func (s *Server) Start() error {
	logger.Info("GRPC server listening at %s", s.Addr())

	lis, err := net.Listen("tcp", s.Addr())
	if err != nil {
		return err
	}

	return s.server.Serve(lis)
}
