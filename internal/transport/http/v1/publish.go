package v1

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/Mort4lis/message-broker/internal/core"
)

func (h *Handler) Publish(w http.ResponseWriter, req *http.Request) {
	ps := httprouter.ParamsFromContext(req.Context())

	queueName := ps.ByName("queue_name")
	logger := h.logger.With(slog.String("queue_name", queueName))

	queue, ok := h.reg.GetByName(queueName)
	if !ok {
		logger.Error("queue is not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	payload, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Error("failed to read body", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = queue.Append(payload); err != nil {
		logger.Error("failed to append message to queue", slog.Any("error", err))
		if errors.Is(err, core.ErrQueueOverflowed) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
