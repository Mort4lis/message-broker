package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/Mort4lis/message-broker/internal/config"
	"github.com/Mort4lis/message-broker/internal/core"
	v1 "github.com/Mort4lis/message-broker/internal/transport/http/v1"
)

type Server struct {
	srv *http.Server
}

func NewServer(conf config.HTTPServer, reg *core.QueueRegistry) *Server {
	h := v1.NewHandler(reg)

	router := httprouter.New()
	router.Handler(http.MethodPost, "/v1/queues/:queue_name/messages", http.HandlerFunc(h.Publish))
	router.Handler(http.MethodPost, "/v1/queues/:queue_name/subscriptions", http.HandlerFunc(h.Subscribe))

	srv := &http.Server{
		Addr:    conf.Listen,
		Handler: router,
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
	return s.srv.Shutdown(ctx)
}
