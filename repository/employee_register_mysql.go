package repository

import (
	"fmt"
	"time"

	"github.com/pikomonde/pikopos/entity"
	log "github.com/sirupsen/logrus"
)

// CreateEmployeeRegister is used to create employee_register
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

// IsEmployeeRegisterExist is used to check whether employee_register exist or not
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
