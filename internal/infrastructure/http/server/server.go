package server

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	server "github.com/antsrp/fio_service/internal/infrastructure/http"
)

type Server struct {
	srv    *http.Server
	logger *zap.Logger
	end    chan bool
	quit   chan os.Signal
}

func NewServer(settings *server.Settings, logger *zap.Logger, end chan bool, h http.Handler) *Server {
	addr := net.JoinHostPort(settings.Host, settings.Port)
	return &Server{
		srv:    &http.Server{Addr: addr, Handler: h},
		logger: logger,
		end:    end,
		quit:   make(chan os.Signal),
	}
}

func (s Server) Start() {

	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(s.quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-s.quit
		s.logger.Info("signal detected, server to shutdown")

		if err := s.srv.Shutdown(context.Background()); err != nil {
			s.logger.Sugar().Fatalf("could not shutdown server: %s", err)
		}
		s.end <- true
	}()

	if err := s.srv.ListenAndServe(); err != nil {
		s.logger.Sugar().Infof("can't listen and serve server: %s", err)
	}
}
