package repository

import (
	"fmt"

	"github.com/pikomonde/pikopos/common"
	"github.com/pikomonde/pikopos/entity"
	log "github.com/sirupsen/logrus"
)

// CreateEmployee is used to create an employee
func (r Repository) CreateEmployee(dbtx common.DBTx, e entity.Employee) (*entity.Employee, error) {
	query := `insert into employee (company_id, full_name, email, phone_number, password, role_id, status) 
	values (?, ?, ?, ?, ?, ?, ?)`
	if dbtx == nil {
		dbtx = r.Clients.PikoposMySQLCli
	}

	res, err := dbtx.Exec(query,
		e.CompanyID, e.FullName, e.Email,
		e.PhoneNumber, "", e.RoleID,
		e.Status.String())
	if err != nil {
		log.WithFields(log.Fields{
			"companyID":      e.CompanyID,
			"employeeStatus": e.Status.String(),
			"employee":       fmt.Sprintf("%+v", e),
		}).Errorln("[Repository][CreateEmployee]: ", err.Error())
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.WithFields(log.Fields{
			"companyID":      e.CompanyID,
			"employeeStatus": e.Status.String(),
			"employee":       fmt.Sprintf("%+v", e),
		}).Errorln("[Repository][CreateEmployee]: ", err.Error())
		return nil, err
	}

	e.ID = int(lastID)
	return &e, nil
}

// GetEmployeeByIdentifier is used to get employee by email or phone number
func (r Repository) GetEmployeeByIdentifier(dbtx common.DBTx, companyID int, employeeIdentifier string) (employee entity.Employee, err error) {
	query := `select company_id, id, full_name, email, phone_number, role_id, status+0
	  from employee where company_id = ? and (email = ? or phone_number = ?)`
	if dbtx == nil {
		dbtx = r.Clients.PikoposMySQLCli
	}

	err = dbtx.QueryRow(query, companyID, employeeIdentifier, employeeIdentifier).Scan(
		&employee.CompanyID, &employee.ID, &employee.FullName,
		&employee.Email, &employee.PhoneNumber, &employee.RoleID,
		&employee.Status,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"companyID":          companyID,
			"employeeIdentifier": employeeIdentifier,
		}).Errorln("[Repository][GetEmployeeByLogin]: ", err.Error())
		return employee, err
	}

	return employee, nil
}

// GetEmployeePassword is used to get employee's hashed password
func (r Repository) GetEmployeePassword(dbtx common.DBTx, companyID, employeeID int) (password string, err error) {
	query := `select password from employee where company_id = ? and id = ?`
	if dbtx == nil {
		dbtx = r.Clients.PikoposMySQLCli
	}

	err = dbtx.QueryRow(query, companyID, employeeID).Scan(&password)
	if err != nil {
		log.WithFields(log.Fields{
			"companyID":  companyID,
			"employeeID": employeeID,
		}).Errorln("[Repository][GetEmployeePassword]: ", err.Error())
		return password, err
	}

	return password, nil
}

// GetEmployeesCount is used to count all employees for pagination
func (r Repository) GetEmployeesCount(dbtx common.DBTx, companyID int) (n int, err error) {
	query := `select count(*) as n from employee where company_id = ?`
	if dbtx == nil {
		dbtx = r.Clients.PikoposMySQLCli
	}

	err = dbtx.QueryRow(query, companyID).Scan(&n)
	if err != nil {
		log.WithFields(log.Fields{
			"companyID": companyID,
		}).Errorln("[Repository][GetEmployeesCount]: ", err.Error())
		return 0, err
	}

	return n, err
}

// GetEmployees is used to get all employee from same company
func (r Repository) GetEmployees(dbtx common.DBTx, companyID int, p Pagination) (employees []entity.Employee, err error) {
	query := `
	  select company_id, id, full_name, email, phone_number, role_id, status+0
		from employee where company_id = ? and id > ? order by id asc limit ?`
	if dbtx == nil {
		dbtx = r.Clients.PikoposMySQLCli
	}

	rows, err := dbtx.Query(query, companyID, p.LastID, p.Limit)
	if err != nil {
		log.WithFields(log.Fields{
			"companyID":  companyID,
			"pagination": fmt.Sprintf("%+v", p),
		}).Errorln("[Repository][GetEmployees]: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		employee := entity.Employee{}
		err = rows.Scan(
			&employee.CompanyID, &employee.ID, &employee.FullName,
			&employee.Email, &employee.PhoneNumber, &employee.RoleID,
			&employee.Status,
		)
		if err != nil {
			log.WithFields(log.Fields{
				"companyID":  companyID,
				"pagination": fmt.Sprintf("%+v", p),
				"count":      len(employees),
			}).Errorln("[Repository][GetEmployees][Scan]: ", err.Error())
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

// UpdateEmployeePassword is used to update employee's hashed password
func (r Repository) UpdateEmployeePassword(dbtx common.DBTx, companyID, employeeID int, password string) error {
	query := `update employee set password = ? where company_id = ? and id = ?`
	if dbtx == nil {
		dbtx = r.Clients.PikoposMySQLCli
	}

	_, err := dbtx.Exec(query, password, companyID, employeeID)
	if err != nil {
		log.WithFields(log.Fields{
			"companyID":  companyID,
			"employeeID": employeeID,
			"password":   password,
		}).Errorln("[Repository][UpdateEmployeePassword]: ", err.Error())
		return err
	}

	return nil
}

// UpdateEmployeeStatus is used to update employee status
func (r Repository) UpdateEmployeeStatus(dbtx common.DBTx, companyID, employeeID int, status int) error {
	query := `update employee set status = ? where company_id = ? and id = ?`
	if dbtx == nil {
		dbtx = r.Clients.PikoposMySQLCli
	}

	_, err := dbtx.Exec(query, status, companyID, employeeID)
	if err != nil {
		log.WithFields(log.Fields{
			"companyID":  companyID,
			"employeeID": employeeID,
			"status":     status,
		}).Errorln("[Repository][UpdateEmployeeStatus]: ", err.Error())
		return err
	}

	return nil
}
