package handler

import (
	"net/http"

	"github.com/pikomonde/pikopos/service"
)

// Handler is used to handles endpoint
type Handler struct {
	Service *service.Service
	Mux     *http.ServeMux
}
