package repository

import (
	sql "github.com/jmoiron/sqlx"
	"github.com/pikomonde/pikopos/common"
	"github.com/pikomonde/pikopos/entity"
	log "github.com/sirupsen/logrus"
)

// CreateRole is used to create a new role
func (r Repository) CreateRole(dbtx common.DBTx, role entity.Role) (*entity.Role, error) {
	query := `insert into role (company_id, name, status) values (?, ?, ?)`
	if dbtx == nil {
		dbtx = r.Clients.PikoposMySQLCli
	}

	res, err := dbtx.Exec(query, role.CompanyID, role.Name, role.Status.String())
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

// func (r Repository) GetRoles(role entity.Role) (*entity.Role, error) {
// 	query := `insert into role (company_id, name, status) values (?, ?, ?)`

// 	res, err := r.Clients.PikoposMySQLCli.Exec(query, role.CompanyID, role.Name, role.Status.String())
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

// GetRolesByIDs is used to get list of roles by role id
func (r Repository) GetRolesByIDs(dbtx common.DBTx, companyID int, ids []int) (map[int]entity.Role, error) {
	query, args, err := sql.In(`select company_id, id, name, status+0
	  from role where company_id = ? and id in (?)`, companyID, ids)
	if err != nil {
		log.WithFields(log.Fields{
			"companyID": companyID,
			"ids":       ids,
		}).Errorln("[Repository][GetRolesByIDs][In]: ", err.Error())
		return nil, err
	}
	query = r.Clients.PikoposMySQLCli.Rebind(query)
	if dbtx == nil {
		dbtx = r.Clients.PikoposMySQLCli
	}

	rows, err := dbtx.Query(query, args...)
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

// GetRoleByID is used to get role by role id
func (r Repository) GetRoleByID(dbtx common.DBTx, companyID int, id int) (role entity.Role, err error) {
	query := `select company_id, id, name, status+0 from role where company_id = ? and id = ?`
	if dbtx == nil {
		dbtx = r.Clients.PikoposMySQLCli
	}

	err = dbtx.QueryRow(query, companyID, id).
		Scan(&role.CompanyID, &role.ID, &role.Name, &role.Status)
	if err != nil {
		log.WithFields(log.Fields{
			"companyID": companyID,
			"id":        id,
		}).Errorln("[Repository][GetRoleByID]: ", err.Error())
		return role, err
	}

	return role, nil
}
