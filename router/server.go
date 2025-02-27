package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/guneyin/app-sdk/logger"
	"github.com/guneyin/app-sdk/utils"
)

const defaultRequestTimeout = time.Second * 5

type Server struct {
	logger     *logger.Logger
	router     *Router
	httpServer *http.Server
}

func New(port string, logger *logger.Logger) *Server {
	mux := http.NewServeMux()
	return &Server{
		logger: logger,
		router: newRouter(logger, mux),
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", utils.GetOutboundIP(), port),
			Handler:      mux,
			ReadTimeout:  defaultRequestTimeout,
			WriteTimeout: defaultRequestTimeout,
		},
	}
}

func (s *Server) SetTimeout(timeout time.Duration) {
	s.httpServer.ReadTimeout = timeout
	s.httpServer.WriteTimeout = timeout
}

func (s *Server) Start() error {
	logger.Info("HTTP server listening at")
	logger.Link("%s", s.Addr())

	return s.httpServer.ListenAndServe()
}

func (s *Server) RegisterHandler(method HTTPMethod, path string, handler HandlerFunc) {
	s.router.registerHandler(method, path, handler)
}

func (s *Server) Addr() string {
	return fmt.Sprintf("http://%s", s.httpServer.Addr)
}
