package v1

import (
	"log/slog"

	"github.com/Mort4lis/message-broker/internal/core"
)

type Handler struct {
	logger *slog.Logger
	reg    *core.QueueRegistry
}

func NewHandler(logger *slog.Logger, reg *core.QueueRegistry) *Handler {
	return &Handler{logger: logger, reg: reg}
}
