package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pikomonde/pikopos/common"
	"github.com/pikomonde/pikopos/entity"
	log "github.com/sirupsen/logrus"
)

const (
	otpExpirationDuration = 60 * time.Minute
)

// RegisterInput is used as parameter for Register. The tag max and min
// is defined from frontend.
type RegisterInput struct {
	CompanyUsername     string `max:"16" json:"company_username"`
	CompanyName         string `max:"32" json:"company_name"`
	FullName            string `max:"32" json:"full_name"`
	EmployeeEmail       string `max:"48" json:"employee_email"`
	EmployeePhoneNumber string `max:"16" json:"employee_phone_number"`
	Password            string `min:"8" json:"password"`
}

// Register is used to register a new user
func (s *Service) Register(ri RegisterInput) (int, error) {
	// TODO: validate input
	// TODO: change to informative error in user
	// TODO: validate email
	// TODO: clean phonenumber (remove "+" in "+62", replace "0" to "62")
	// TODO: add validation, whether a company_username already taken or not
	// remove password for security while logging
	passwordRaw := ri.Password
	ri.Password = ""
	tx, err := s.Repository.Clients.PikoposMySQLCli.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"registerInput": fmt.Sprintf("%+v", ri),
		}).Errorln("[Service][Register][Begin]: ", err.Error())
		return http.StatusInternalServerError, err
	}

	company, err := s.Repository.CreateCompany(tx, entity.Company{
		Username: ri.CompanyUsername,
		Name:     ri.CompanyName,
		Status:   entity.CompanyStatusFree,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"registerInput": fmt.Sprintf("%+v", ri),
			"rollback":      tx.Rollback(),
		}).Errorln("[Service][Register][CreateCompany]: ", err.Error())
		return http.StatusInternalServerError, err
	}

	role, err := s.Repository.CreateRole(tx, entity.Role{
		CompanyID: company.ID,
		Name:      entity.RoleSuperAdmin,
		Status:    entity.RoleStatusActive,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"registerInput": fmt.Sprintf("%+v", ri),
			"company":       fmt.Sprintf("%+v", company),
			"rollback":      tx.Rollback(),
		}).Errorln("[Service][Register][CreateRole]: ", err.Error())
		return http.StatusInternalServerError, err
	}

	employee, err := s.Repository.CreateEmployee(tx, entity.Employee{
		CompanyID:   company.ID,
		FullName:    ri.FullName,
		Email:       ri.EmployeeEmail,
		PhoneNumber: ri.EmployeePhoneNumber,
		RoleID:      role.ID,
		Status:      entity.EmployeeStatusUnverified,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"registerInput": fmt.Sprintf("%+v", ri),
			"company":       fmt.Sprintf("%+v", company),
			"role":          fmt.Sprintf("%+v", role),
			"rollback":      tx.Rollback(),
		}).Errorln("[Service][Register][CreateEmployee]: ", err.Error())
		return http.StatusInternalServerError, err
	}

	passwordHashed := common.SHA256(fmt.Sprintf("%s-%s-%s-%d",
		passwordRaw, ri.EmployeeEmail, ri.EmployeePhoneNumber, employee.ID))
	err = s.Repository.UpdateEmployeePassword(tx, company.ID, employee.ID, passwordHashed)
	if err != nil {
		log.WithFields(log.Fields{
			"registerInput": fmt.Sprintf("%+v", ri),
			"company":       fmt.Sprintf("%+v", company),
			"role":          fmt.Sprintf("%+v", role),
			"employee":      fmt.Sprintf("%+v", employee),
			"rollback":      tx.Rollback(),
		}).Errorln("[Service][Register][UpdateEmployeePassword]: ", err.Error())
		return http.StatusInternalServerError, err
	}

	otpCode, err := common.OTP(12)
	if err != nil {
		log.WithFields(log.Fields{
			"registerInput": fmt.Sprintf("%+v", ri),
			"company":       fmt.Sprintf("%+v", company),
			"role":          fmt.Sprintf("%+v", role),
			"employee":      fmt.Sprintf("%+v", employee),
		}).Errorln("[Service][Register][OTP6]emailOTP: ", err.Error())
		return 0, err
	}
	// TODO: sent to email and phone number
	log.Debugln("[Register] OTP:", otpCode)

	otpHashed := common.SHA256(fmt.Sprintf("%s-%s-%s-%d",
		otpCode, ri.EmployeeEmail, ri.EmployeePhoneNumber, employee.ID))
	_, err = s.Repository.CreateEmployeeRegister(tx, entity.EmployeeRegister{
		EmployeeID:  employee.ID,
		Email:       ri.EmployeeEmail,
		PhoneNumber: ri.EmployeePhoneNumber,
		OTPCode:     otpHashed,
		ExpiredAt:   time.Now().Add(otpExpirationDuration),
	})
	if err != nil {
		log.WithFields(log.Fields{
			"registerInput": fmt.Sprintf("%+v", ri),
			"company":       fmt.Sprintf("%+v", company),
			"role":          fmt.Sprintf("%+v", role),
			"employee":      fmt.Sprintf("%+v", employee),
			"rollback":      tx.Rollback(),
		}).Errorln("[Service][Register][UpdateEmployeePassword]: ", err.Error())
		return http.StatusInternalServerError, err
	}

	err = tx.Commit()
	if err != nil {
		log.WithFields(log.Fields{
			"registerInput": fmt.Sprintf("%+v", ri),
			"company":       fmt.Sprintf("%+v", company),
			"role":          fmt.Sprintf("%+v", role),
			"employee":      fmt.Sprintf("%+v", employee),
		}).Errorln("[Service][Register][Commit]: ", err.Error())
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// VerifyInput is used as input for user in registration process
type VerifyInput struct {
	EmployeeID     int    `json:"employee_id"`
	OTPEmail       string `min:"6" max:"6" json:"otp_email"`
	OTPPhoneNumber string `min:"6" max:"6" json:"otp_phone_number"`
}

// Verify is used to verify a newly registerd user
func (s *Service) Verify(vi VerifyInput) (int, error) {
	// TODO: validate input
	// TODO: change to informative error in user
	// TODO: use otpHashed instead of otpCode

	otpCode := vi.OTPEmail + vi.OTPPhoneNumber
	// otpHashed := common.SHA256(fmt.Sprintf("%s-%s-%s-%d",
	// 	otpCode, ri.EmployeeEmail, ri.PhoneNumber, employee.ID))
	isEmployeeRegisterExist, err := s.Repository.IsEmployeeRegisterExist(nil, vi.EmployeeID, otpCode)
	if err != nil {
		log.WithFields(log.Fields{
			"verifyInput": fmt.Sprintf("%+v", vi),
		}).Errorln("[Service][Verify][IsEmployeeRegisterExist]: ", err.Error())
		return http.StatusInternalServerError, err
	}
	if !isEmployeeRegisterExist {
		log.WithFields(log.Fields{
			"verifyInput": fmt.Sprintf("%+v", vi),
		}).Errorln("[Service][Verify][IsEmployeeRegisterExist]: wrong otp code")
		return http.StatusBadRequest, err
	}

	// TODO: change employee status
	// s.Repository.UpdateEmployeeStatus()

	return http.StatusOK, nil
}
