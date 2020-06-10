package repository

import (
	// initialize mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/pikomonde/pikopos/clients"
	"github.com/pikomonde/pikopos/common"
	"github.com/pikomonde/pikopos/entity"
	"github.com/pikomonde/pikopos/repository/localconfig"
	"github.com/pikomonde/pikopos/repository/mysqlredis"
	"golang.org/x/oauth2"
)

type RepositoryAuth interface {
	GetAuthConfig(provider string) (*oauth2.Config, error)
}

type RepositoryCompany interface {
	CreateCompany(dbtx common.DBTx, company entity.Company) (*entity.Company, error)
	GetCompanyByUsername(dbtx common.DBTx, companyUsername string) (company entity.Company, err error)
}

type RepositoryEmployee interface {
	CreateEmployee(dbtx common.DBTx, e entity.Employee) (*entity.Employee, error)
	GetEmployeeByIdentifier(dbtx common.DBTx, companyID int, employeeIdentifier string) (employee entity.Employee, err error)
	IsEmployeeExist(dbtx common.DBTx, companyID int, email, phoneNumber string) (bool, error)
	GetEmployeePassword(dbtx common.DBTx, companyID, employeeID int) (password string, err error)
	GetEmployeesCount(dbtx common.DBTx, companyID int) (n int, err error)
	GetEmployees(dbtx common.DBTx, companyID int, p common.PaginationRepo) (employees []entity.Employee, err error)
	UpdateEmployee(dbtx common.DBTx, e entity.Employee) (int, *entity.Employee, error)
	UpdateEmployeePassword(dbtx common.DBTx, companyID, employeeID int, password string) error
	UpdateEmployeeStatus(dbtx common.DBTx, companyID, employeeID int, status int) error
}

type RepositoryRole interface {
	CreateRole(dbtx common.DBTx, role entity.Role) (*entity.Role, error)
	// GetRoles(role entity.Role) (*entity.Role, error)
	GetRolesByIDs(dbtx common.DBTx, companyID int, ids []int) (map[int]entity.Role, error)
	GetRoleByID(dbtx common.DBTx, companyID int, id int) (role entity.Role, err error)
}

// === MySQLRedis repositories ===

// NewLocalConfigAuth returns Auth repository from config
func NewLocalConfigAuth(c *clients.Clients) RepositoryAuth {
	return &localconfig.RepositoryAuth{Clients: c}
}

// NewMySQLRedisCompany returns Company repository using mysql connection
func NewMySQLRedisCompany(c *clients.Clients) RepositoryCompany {
	return &mysqlredis.RepositoryCompany{Clients: c}
}

// NewMySQLRedisEmployee returns Employee repository using mysql connection
func NewMySQLRedisEmployee(c *clients.Clients) RepositoryEmployee {
	return &mysqlredis.RepositoryEmployee{Clients: c}
}

// NewMySQLRedisRole returns Role repository using mysql connection
func NewMySQLRedisRole(c *clients.Clients) RepositoryRole {
	return &mysqlredis.RepositoryRole{Clients: c}
}
