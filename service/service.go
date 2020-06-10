package service

import (
	// initialize mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/pikomonde/pikopos/repository"
	sAuth "github.com/pikomonde/pikopos/service/auth"
	sEmployee "github.com/pikomonde/pikopos/service/employee"
)

// NewAuth returns the AuthRegister service
func NewAuth(
	rAuth repository.RepositoryAuth,
	rCompany repository.RepositoryCompany,
	rEmployee repository.RepositoryEmployee,
	rRole repository.RepositoryRole,
) *sAuth.ServiceAuth {
	return &sAuth.ServiceAuth{
		RepositoryAuth:     rAuth,
		RepositoryCompany:  rCompany,
		RepositoryEmployee: rEmployee,
		RepositoryRole:     rRole,
	}
}

// NewEmployee returns the Employee service
func NewEmployee(
	rEmployee repository.RepositoryEmployee,
	rRole repository.RepositoryRole,
) *sEmployee.ServiceEmployee {
	return &sEmployee.ServiceEmployee{
		RepositoryEmployee: rEmployee,
		RepositoryRole:     rRole,
	}
}
