package employee

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pikomonde/pikopos/common"
	"github.com/pikomonde/pikopos/entity"
	log "github.com/sirupsen/logrus"
)

// EmployeeCreateInput is used as request for create new employee
type EmployeeCreateInput struct {
	CompanyID   int    `json:"-"`
	FullName    string `max:"32" json:"full_name"`
	Email       string `max:"48" json:"email"`
	PhoneNumber string `max:"16" json:"phone_number"`
	RoleID      int    `json:"role_id"`
	// OutletIDs   []int  `json:"outlet_ids"`
}

// EmployeeUpdateInput is used as request for update existing employee
type EmployeeUpdateInput struct {
	CompanyID int    `json:"-"`
	ID        int    `json:"id"`
	FullName  string `max:"32" json:"full_name"`
	RoleID    int    `json:"role_id"`
	// OutletIDs   []int  `json:"outlet_ids"`
}

// EmployeeOutput is used as response for employee
type EmployeeOutput struct {
	ID          int    `json:"id"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	RoleID      int    `json:"role_id"`
	RoleName    string `json:"role_name"`
	Status      string `json:"status"`
	RowUpdated  int    `json:"row_updated"`
}

// CreateEmployee is used to create employee data
func (s *ServiceEmployee) CreateEmployee(eci EmployeeCreateInput) (*EmployeeOutput, int, error) {
	// TODO: validate input
	// TODO: change to informative error in user
	// TODO: prevent creating user with same email or phone_number in a company

	isExist, err := s.repositoryEmployee.IsEmployeeExist(nil, eci.CompanyID, eci.Email, eci.PhoneNumber)
	if err != nil {
		log.WithFields(log.Fields{
			"employeeInput": fmt.Sprintf("%+v", eci),
		}).Errorln("[ServiceEmployee][CreateEmployee][IsEmployeeExist]: ", err.Error())
		return nil, http.StatusInternalServerError, err
	}
	if isExist {
		return nil, http.StatusBadRequest, errors.New(common.ErrorIdentifierAlreadyExist)
	}

	employee, err := s.repositoryEmployee.CreateEmployee(nil, entity.Employee{
		CompanyID:   eci.CompanyID,
		FullName:    eci.FullName,
		Email:       eci.Email,
		PhoneNumber: eci.PhoneNumber,
		RoleID:      eci.RoleID,
		Status:      entity.EmployeeStatusUnverified,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"employeeInput": fmt.Sprintf("%+v", eci),
			"isExist":       isExist,
		}).Errorln("[ServiceEmployee][CreateEmployee][CreateEmployee]: ", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	role, err := s.repositoryRole.GetRoleByID(nil, employee.CompanyID, employee.RoleID)
	if err != nil {
		log.WithFields(log.Fields{
			"employeeInput": fmt.Sprintf("%+v", eci),
			"isExist":       isExist,
			"employee":      fmt.Sprintf("%+v", employee),
		}).Errorln("[ServiceEmployee][CreateEmployee][GetRoleByID]: ", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return &EmployeeOutput{
		ID:          employee.ID,
		FullName:    employee.FullName,
		Email:       employee.Email,
		PhoneNumber: employee.PhoneNumber,
		RoleID:      employee.RoleID,
		RoleName:    role.Name,
		Status:      employee.Status.String(),
	}, http.StatusOK, nil
}

// UpdateEmployee is used to update employee data
func (s *ServiceEmployee) UpdateEmployee(eui EmployeeUpdateInput) (*EmployeeOutput, int, error) {
	// TODO: validate input
	// TODO: change to informative error in user

	cnt, employee, err := s.repositoryEmployee.UpdateEmployee(nil, entity.Employee{
		CompanyID: eui.CompanyID,
		ID:        eui.ID,
		FullName:  eui.FullName,
		RoleID:    eui.RoleID,
		Status:    entity.EmployeeStatusUnverified,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"employeeInput": fmt.Sprintf("%+v", eui),
		}).Errorln("[ServiceEmployee][UpdateEmployee][UpdateEmployee]: ", err.Error())
		return nil, http.StatusInternalServerError, err
	}
	// if cnt == 0 {
	// 	return nil, http.StatusBadRequest, errors.New(errorNotUpdateRowNotExist)
	// }

	role, err := s.repositoryRole.GetRoleByID(nil, employee.CompanyID, employee.RoleID)
	if err != nil {
		log.WithFields(log.Fields{
			"employeeInput": fmt.Sprintf("%+v", eui),
			"employee":      fmt.Sprintf("%+v", employee),
		}).Errorln("[ServiceEmployee][UpdateEmployee][GetRoleByID]: ", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return &EmployeeOutput{
		ID:          employee.ID,
		FullName:    employee.FullName,
		Email:       employee.Email,
		PhoneNumber: employee.PhoneNumber,
		RoleID:      employee.RoleID,
		RoleName:    role.Name,
		Status:      employee.Status.String(),
		RowUpdated:  cnt,
	}, http.StatusOK, nil
}
