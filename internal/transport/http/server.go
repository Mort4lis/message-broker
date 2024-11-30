package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/Mort4lis/message-broker/internal/config"
	"github.com/Mort4lis/message-broker/internal/core"
	v1 "github.com/Mort4lis/message-broker/internal/transport/http/v1"
	"github.com/Mort4lis/message-broker/pkg/httputils/middleware"
)

type Server struct {
	srv *http.Server
}

func NewServer(logger *slog.Logger, conf config.HTTPServer, reg *core.QueueRegistry) *Server {
	h := v1.NewHandler(logger, reg)

	router := httprouter.New()
	router.Handler(http.MethodPost, "/v1/queues/:queue_name/messages", http.HandlerFunc(h.Publish))
	router.Handler(http.MethodPost, "/v1/queues/:queue_name/subscriptions", http.HandlerFunc(h.Subscribe))

	srv := &http.Server{
		Addr: conf.Listen,
		Handler: middleware.Timeout(
			conf.RequestTimeout,
			middleware.Log(logger, router),
		),
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
		ReadTimeout:       conf.ReadTimeout,
		ReadHeaderTimeout: conf.ReadHeaderTimeout,
		WriteTimeout:      conf.WriteTimeout,
	}
	return &Server{srv: srv}
}

func (s *Server) Run() error {
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %v", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx) //nolint:wrapcheck // no need to wrap error
}
