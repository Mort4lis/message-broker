package v1

import (
	"github.com/Mort4lis/message-broker/internal/core"
)

type Handler struct {
	reg *core.QueueRegistry
}

func NewHandler(reg *core.QueueRegistry) *Handler {
	return &Handler{reg: reg}
}
