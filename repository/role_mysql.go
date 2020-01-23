package repository

import (
	"github.com/pikomonde/pikopos/entity"
	log "github.com/sirupsen/logrus"
)

func (c Repository) CreateRole(role entity.Role) (*entity.Role, error) {
	query := `insert into role (company_id, name, status) values (?, ?, ?)`

	res, err := c.Clients.PikoposMySQLCli.Exec(query, role.CompanyID, role.Name, role.Status.String())
	if err != nil {
		log.WithFields(log.Fields{
			"companyID":  role.CompanyID,
			"roleName":   role.Name,
			"roleStatus": role.Status.String(),
		}).Errorln("[Repository][CreateRole]: ", err.Error())
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.WithFields(log.Fields{
			"companyID":  role.CompanyID,
			"roleName":   role.Name,
			"roleStatus": role.Status.String(),
		}).Errorln("[Repository][CreateRole][LastInsertId]: ", err.Error())
		return nil, err
	}

	role.ID = int(lastID)
	return &role, nil
}
