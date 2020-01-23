package repository

import (
	// initialize mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/pikomonde/pikopos/clients"
)

// Repository contains clients and all repositories
type Repository struct {
	Clients *clients.Clients
}

// New returns the repository
func New(c *clients.Clients) *Repository {
	return &Repository{c}
}
