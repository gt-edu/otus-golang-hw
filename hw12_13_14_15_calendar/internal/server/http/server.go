package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/gt-edu/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
)

type AppHandler struct {
	logger logger.Logger
}

func (h *AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("EVERYTHING IS FINE\n"))
	if err != nil {
		h.logger.Error(err.Error())
		return
	}
}

type Server struct { // TODO
	httpServer *http.Server
	config     *config.Config
	logger     logger.Logger
}

type Application interface { // TODO
}

func NewServer(appLogger logger.Logger, app Application, appConfig *config.Config) *Server {
	return &Server{
		logger: appLogger,
		config: appConfig,
	}
}

func (s *Server) Start() error {
	handler := &AppHandler{
		logger: s.logger,
	}
	addr := net.JoinHostPort(s.config.HTTP.Hostname, s.config.HTTP.Port)
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      loggingMiddleware(handler, s.logger),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.logger.Info(fmt.Sprintf("Starting a server on %s:%s", s.config.HTTP.Hostname, s.config.HTTP.Port))

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
