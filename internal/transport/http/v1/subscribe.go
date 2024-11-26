package v1

import (
	"net/http"
)

func (h *Handler) Subscribe(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(req.URL.String()))
}
