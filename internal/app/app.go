package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/Mort4lis/message-broker/internal/config"
	"github.com/Mort4lis/message-broker/internal/core"
	"github.com/Mort4lis/message-broker/internal/logging"
	"github.com/Mort4lis/message-broker/internal/transport/http"
)

func Run(confPath string) error {
	var conf config.Config
	if err := cleanenv.ReadConfig(confPath, &conf); err != nil {
		return fmt.Errorf("read config: %v", err)
	}

	logger, err := logging.NewLoggerFromConfig(conf.Logging)
	if err != nil {
		return fmt.Errorf("create logger: %v", err)
	}

	registry := core.NewQueueRegistry()
	for _, queueConf := range conf.Queues {
		queue := core.NewQueue(queueConf.Name, queueConf.MaxMessages, queueConf.MaxSubscribers)
		if err = registry.Register(queue); err != nil {
			return fmt.Errorf("register queue with name %s: %v", queue.Name(), err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	registry.ForEach(func(q *core.Queue) bool {
		q.StartConsume(ctx)
		return true
	})

	errCh := make(chan error)
	httpServer := http.NewServer(logger, conf.HTTPServer, registry)
	go func() {
		logger.Info("Listen http server", slog.String("addr", conf.HTTPServer.Listen))
		if err = httpServer.Run(); err != nil {
			errCh <- fmt.Errorf("run http server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)

	select {
	case err = <-errCh:
		return err
	case sig := <-quit:
		logger.Info("Caught signal. Shutting down...", slog.String("signal", sig.String()))
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), conf.ShutdownTimeout)
		defer shutdownCancel()

		if err = httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("Failed to shutdown http server", slog.String("error", err.Error()))
		}

		cancel()
		registry.ForEach(func(q *core.Queue) bool {
			q.Wait()
			return true
		})
	}

	return nil
}
