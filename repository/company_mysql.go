package repository

import (
	"github.com/pikomonde/pikopos/entity"
	log "github.com/sirupsen/logrus"
)

func (c Repository) CreateCompany(company entity.Company) (*entity.Company, error) {
	query := `insert into company (username, name, status) values (?, ?, ?)`

	res, err := c.Clients.PikoposMySQLCli.Exec(query, company.Username, company.Name, company.Status.String())
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

func (c Repository) GetCompanyByUsername(companyUsername string) (company entity.Company, err error) {
	query := `select id, username, name, status+0 from company where username = ?`

	err = c.Clients.PikoposMySQLCli.QueryRow(query, companyUsername).
		Scan(&company.ID, &company.Username, &company.Name, &company.Status)
	if err != nil {
		log.WithFields(log.Fields{
			"companyUsername": companyUsername,
		}).Errorln("[Repository][GetCompanyByID]: ", err.Error())
		return company, err
	}

	return company, nil
}
