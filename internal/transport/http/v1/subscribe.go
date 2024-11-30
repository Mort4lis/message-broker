package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/Mort4lis/message-broker/internal/core"
)

func (h *Handler) Subscribe(w http.ResponseWriter, req *http.Request) {
	ps := httprouter.ParamsFromContext(req.Context())

	queueName := ps.ByName("queue_name")
	logger := h.logger.With(slog.String("queue_name", queueName))

	queue, ok := h.reg.GetByName(queueName)
	if !ok {
		logger.Error("queue is not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cons, err := queue.Subscribe()
	if errors.Is(err, core.ErrReachedSubscriberLimit) {
		logger.Error("reached subscriber limit")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		logger.Error("failed to subscribe on queue", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer cons.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	isFirst := true

	for {
		msg, readErr := cons.ReadMessage(req.Context())
		if readErr != nil {
			if !errors.Is(readErr, req.Context().Err()) {
				logger.Error("failed to read message", slog.Any("error", readErr))
			}

			return
		}

		if !isFirst {
			w.Write([]byte("\n"))
		}

		w.Write(msg)
		isFirst = false
	}
}
