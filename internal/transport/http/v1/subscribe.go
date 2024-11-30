package v1

import (
	"net/http"
)

func (h *Handler) Subscribe(w http.ResponseWriter, req *http.Request) {
	_, _ = w, req
}
