package repository

import (
	"github.com/pikomonde/pikopos/common"
	"github.com/pikomonde/pikopos/entity"
	log "github.com/sirupsen/logrus"
)

// CreateCompany is used to create a company when user first register
func (r Repository) CreateCompany(dbtx common.DBTx, company entity.Company) (*entity.Company, error) {
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
		}).Errorln("[Repository][CreateCompany]: ", err.Error())
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.WithFields(log.Fields{
			"companyUsername": company.Username,
			"companyName":     company.Name,
			"companyStatus":   company.Status.String(),
		}).Errorln("[Repository][CreateCompany][LastInsertId]: ", err.Error())
		return nil, err
	}

	company.ID = int(lastID)
	return &company, nil
}

// GetCompanyByUsername is used to get company by company username
func (r Repository) GetCompanyByUsername(dbtx common.DBTx, companyUsername string) (company entity.Company, err error) {
	query := `select id, username, name, status+0 from company where username = ?`
	if dbtx == nil {
		dbtx = r.Clients.PikoposMySQLCli
	}

	err = dbtx.QueryRow(query, companyUsername).
		Scan(&company.ID, &company.Username, &company.Name, &company.Status)
	if err != nil {
		log.WithFields(log.Fields{
			"companyUsername": companyUsername,
		}).Errorln("[Repository][GetCompanyByID]: ", err.Error())
		return company, err
	}

	return company, nil
}
