package handler

import (
	"net/http"

	sAuth "github.com/pikomonde/pikopos/service/auth"
	sEmployee "github.com/pikomonde/pikopos/service/employee"
)

// Handler is used to handles endpoint
type Handler struct {
	ServiceAuth     *sAuth.ServiceAuth
	ServiceEmployee *sEmployee.ServiceEmployee
	Mux             *http.ServeMux
}
