package delivery

import (
	// initialize mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/pikomonde/pikopos/service"
)

// Delivery contains services and endpoints
type Delivery struct {
	Service      *service.Service
	EchoDelivery *echo.Echo
}

// New returns the delivery
func New(s *service.Service) *Delivery {
	dlvrEcho := echo.New()
	return &Delivery{s, dlvrEcho}
}

// Start starts the delivery server
func (d *Delivery) Start() {
	d.PrepareAuth()
	d.EchoDelivery.Start(":1235")
}
