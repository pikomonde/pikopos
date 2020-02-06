package delivery

import (
	// initialize mysql driver
	"crypto/tls"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pikomonde/pikopos/delivery/handler"
	"github.com/pikomonde/pikopos/service"
)

// Delivery contains services and endpoints
type Delivery struct {
	Handler *handler.Handler
}

// New returns the delivery
func New(s *service.Service) *Delivery {
	mux := http.NewServeMux()
	return &Delivery{&handler.Handler{
		Service: s,
		Mux:     mux,
	}}
}

// Start starts the delivery server
func (d *Delivery) Start() {
	// Register handlers
	d.Handler.RegisterAuth()

	// Starting server
	srv := &http.Server{
		ReadTimeout:  5000 * time.Millisecond,
		WriteTimeout: 5000 * time.Millisecond,
		TLSConfig: &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},
		},
		Handler: d.Handler.Mux,
		Addr:    ":1235",
	}
	srv.ListenAndServe()
	// log.Fatal(srv.ListenAndServe())
}
