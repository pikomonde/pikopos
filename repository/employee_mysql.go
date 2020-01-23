package repository

import (
	"fmt"
	"time"

	"github.com/pikomonde/pikopos/entity"
	log "github.com/sirupsen/logrus"
)

func (c Repository) CreateEmployee(e entity.Employee) (*entity.Employee, error) {
	query := `insert into employee (company_id, first_name, last_name, email, phone_number, password, role_id, status) 
	values (?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := c.Clients.PikoposMySQLCli.Exec(query,
		e.CompanyID, e.FirstName, e.LastName,
		e.Email, e.PhoneNumber, "",
		e.RoleID, e.Status.String())
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

func (c Repository) GetEmployeeByIdentifier(companyID int, employeeIdentifier string) (employee entity.Employee, err error) {
	query := `select company_id, id, first_name, last_name, email, phone_number, role_id, status+0
	  from employee where company_id = ? and (email = ? or phone_number = ?)`

	err = c.Clients.PikoposMySQLCli.QueryRow(query, companyID, employeeIdentifier, employeeIdentifier).Scan(
		&employee.CompanyID, &employee.ID, &employee.FirstName,
		&employee.LastName, &employee.Email, &employee.PhoneNumber,
		&employee.RoleID, &employee.Status,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"employeeIdentifier": employeeIdentifier,
		}).Errorln("[Repository][GetEmployeeByLogin]: ", err.Error())
		return employee, err
	}

	return employee, nil
}

func (c Repository) GetEmployeePassword(companyID, employeeID int) (password string, err error) {
	query := `select password from employee where company_id = ? and id = ?`

	err = c.Clients.PikoposMySQLCli.QueryRow(query, companyID, employeeID).Scan(&password)
	if err != nil {
		log.WithFields(log.Fields{
			"companyID":  companyID,
			"employeeID": employeeID,
		}).Errorln("[Repository][GetEmployeePassword]: ", err.Error())
		return password, err
	}

	return password, nil
}

func (c Repository) UpdateEmployeePassword(companyID, employeeID int, password string) error {
	query := `update employee set password = ? where company_id = ? and id = ?`

	_, err := c.Clients.PikoposMySQLCli.Exec(query, password, companyID, employeeID)
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

func (c Repository) UpdateEmployeeStatus(companyID, employeeID int, status int) error {
	query := `update employee set status = ? where company_id = ? and id = ?`

	_, err := c.Clients.PikoposMySQLCli.Exec(query, status, companyID, employeeID)
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

func (c Repository) CreateEmployeeRegister(er entity.EmployeeRegister) (*entity.EmployeeRegister, error) {
	query := `insert into employee_register (employee_id, email, phone_number, otp_code, expired_at) 
	values (?, ?, ?, ?, ?)`

	res, err := c.Clients.PikoposMySQLCli.Exec(query,
		er.EmployeeID, er.Email, er.PhoneNumber,
		er.OTPCode, er.ExpiredAt)
	if err != nil {
		log.WithFields(log.Fields{
			"employeeID":       er.EmployeeID,
			"employeeRegister": fmt.Sprintf("%+v", er),
		}).Errorln("[Repository][CreateEmployeeRegister]: ", err.Error())
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.WithFields(log.Fields{
			"employeeID":       er.EmployeeID,
			"employeeRegister": fmt.Sprintf("%+v", er),
		}).Errorln("[Repository][CreateEmployeeRegister]: ", err.Error())
		return nil, err
	}

	er.ID = int(lastID)
	return &er, nil
}

func (c Repository) IsEmployeeRegisterExist(employeeID int, otpCode string) (bool, error) {
	query := `select 1 from employee_register
	  where employee_id = ? and otp_code = ? and expired_at < ?`

	isEmployeeRegisterExist := false
	timeNow := time.Now()

	err := c.Clients.PikoposMySQLCli.QueryRow(query, employeeID, otpCode, timeNow).
		Scan(&isEmployeeRegisterExist)
	if err != nil {
		log.WithFields(log.Fields{
			"employeeID": employeeID,
			"time":       timeNow,
		}).Errorln("[Repository][IsEmployeeRegisterExist]: ", err.Error())
		return false, err
	}

	return isEmployeeRegisterExist, nil
}
