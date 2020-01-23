package service

import (
	// initialize mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/pikomonde/pikopos/repository"
)

// Service contains repositories and all use cases
type Service struct {
	Repository *repository.Repository
}

// New returns the service
func New(r *repository.Repository) *Service {
	return &Service{r}
}
