package delivery

import (
	// initialize mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	echoMiddle "github.com/labstack/echo/middleware"
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
	// TODO; remove CORS in production
	dlvrEcho.Use(echoMiddle.CORS())
	return &Delivery{s, dlvrEcho}
}

// Start starts the delivery server
func (d *Delivery) Start() {
	d.PrepareAuth()
	d.EchoDelivery.Start(":1235")
}
