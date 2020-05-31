package mysqlredis

import (
	"github.com/pikomonde/pikopos/clients"
	"github.com/pikomonde/pikopos/common"
	"github.com/pikomonde/pikopos/entity"
	log "github.com/sirupsen/logrus"
)

// RepositoryCompany contains clients and Company repositories
type RepositoryCompany struct {
	Clients *clients.Clients
}

// CreateCompany is used to create a company when user first register
func (r RepositoryCompany) CreateCompany(dbtx common.DBTx, company entity.Company) (*entity.Company, error) {
	query := `insert into company (username, name, status) values (?, ?, ?)`
	if dbtx == nil {
		dbtx = r.Clients.PikoposMySQLCli
	}

	res, err := dbtx.Exec(query, company.Username, company.Name, company.Status.String())
	if err != nil {
		log.WithFields(log.Fields{
			"companyUsername": company.Username,
			"companyName":     company.Name,
			"companyStatus":   company.Status.String(),
		}).Errorln("[RepositoryCompany][CreateCompany]: ", err.Error())
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.WithFields(log.Fields{
			"companyUsername": company.Username,
			"companyName":     company.Name,
			"companyStatus":   company.Status.String(),
		}).Errorln("[RepositoryCompany][CreateCompany][LastInsertId]: ", err.Error())
		return nil, err
	}

	company.ID = int(lastID)
	return &company, nil
}

// GetCompanyByUsername is used to get company by company username
func (r RepositoryCompany) GetCompanyByUsername(dbtx common.DBTx, companyUsername string) (company entity.Company, err error) {
	query := `select id, username, name, status-1 from company where username = ?`
	if dbtx == nil {
		dbtx = r.Clients.PikoposMySQLCli
	}

	err = dbtx.QueryRow(query, companyUsername).
		Scan(&company.ID, &company.Username, &company.Name, &company.Status)
	if err != nil {
		log.WithFields(log.Fields{
			"companyUsername": companyUsername,
		}).Errorln("[RepositoryCompany][GetCompanyByID]: ", err.Error())
		return company, err
	}

	return company, nil
}
