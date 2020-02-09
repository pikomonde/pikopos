package repository

import (
	sql "github.com/jmoiron/sqlx"
	"github.com/pikomonde/pikopos/entity"
	log "github.com/sirupsen/logrus"
)

// CreateRole is used to create a new role
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

// func (c Repository) GetRoles(role entity.Role) (*entity.Role, error) {
// 	query := `insert into role (company_id, name, status) values (?, ?, ?)`

// 	res, err := c.Clients.PikoposMySQLCli.Exec(query, role.CompanyID, role.Name, role.Status.String())
// 	if err != nil {
// 		log.WithFields(log.Fields{
// 			"companyID":  role.CompanyID,
// 			"roleName":   role.Name,
// 			"roleStatus": role.Status.String(),
// 		}).Errorln("[Repository][CreateRole]: ", err.Error())
// 		return nil, err
// 	}

// 	lastID, err := res.LastInsertId()
// 	if err != nil {
// 		log.WithFields(log.Fields{
// 			"companyID":  role.CompanyID,
// 			"roleName":   role.Name,
// 			"roleStatus": role.Status.String(),
// 		}).Errorln("[Repository][CreateRole][LastInsertId]: ", err.Error())
// 		return nil, err
// 	}

// 	role.ID = int(lastID)
// 	return &role, nil
// }

// GetRolesByIDs is used to get list of roles by role
func (c Repository) GetRolesByIDs(companyID int, ids []int) (map[int]entity.Role, error) {
	query, args, err := sql.In(`select company_id, id, name, status+0
	  from role where company_id = ? and id in (?)`, companyID, ids)
	if err != nil {
		log.WithFields(log.Fields{
			"companyID": companyID,
			"ids":       ids,
		}).Errorln("[Repository][GetRolesByIDs][In]: ", err.Error())
		return nil, err
	}

	query = c.Clients.PikoposMySQLCli.Rebind(query)

	rows, err := c.Clients.PikoposMySQLCli.Query(query, args...)
	if err != nil {
		log.WithFields(log.Fields{
			"companyID": companyID,
			"ids":       ids,
			"args":      args,
		}).Errorln("[Repository][GetRolesByIDs]: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	roles := make(map[int]entity.Role, 0)
	for rows.Next() {
		role := entity.Role{}
		err = rows.Scan(&role.CompanyID, &role.ID, &role.Name, &role.Status)
		if err != nil {
			log.WithFields(log.Fields{
				"companyID": companyID,
				"ids":       ids,
				"count":     len(roles),
			}).Errorln("[Repository][GetRolesByIDs][Scan]: ", err.Error())
			return nil, err
		}
		roles[role.ID] = role
	}

	return roles, nil
}
