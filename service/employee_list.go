package service

import (
	"fmt"
	"net/http"

	"github.com/pikomonde/pikopos/repository"
	log "github.com/sirupsen/logrus"
)

// EmployeeListInput is used as request for employee list
type EmployeeListInput struct {
	CompanyID int
	LastID    int `min:"0" json:"last_id"`
	Limit     int `min:"10" json:"limit"`
}

// EmployeeListOutput is used as response for employee list
type EmployeeListOutput struct {
	Count     int              `json:"count"`
	Employees []EmployeeOutput `json:"employees"`
}

// GetEmployeeList is used to get employee list a new user
func (s *Service) GetEmployeeList(eli EmployeeListInput) (*EmployeeListOutput, int, error) {
	// TODO: validate input
	// TODO: change to informative error in user
	tx, err := s.Repository.Clients.PikoposMySQLCli.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"employeeListInput": fmt.Sprintf("%+v", eli),
		}).Errorln("[Service][GetEmployeeList][Begin]: ", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	count, err := s.Repository.GetEmployeesCount(tx, eli.CompanyID)
	if err != nil {
		log.WithFields(log.Fields{
			"employeeListInput": fmt.Sprintf("%+v", eli),
		}).Errorln("[Service][GetEmployeeList][GetEmployees]: ", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	employees, err := s.Repository.GetEmployees(tx, eli.CompanyID, repository.Pagination{
		LastID: eli.LastID,
		Limit:  eli.Limit,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"employeeListInput": fmt.Sprintf("%+v", eli),
		}).Errorln("[Service][GetEmployeeList][GetEmployees]: ", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	// remove roleIDs duplicates
	encountered := make(map[int]bool, 0)
	roleIDs := make([]int, 0)
	for _, v := range employees {
		if !encountered[v.RoleID] {
			encountered[v.RoleID] = true
			roleIDs = append(roleIDs, v.RoleID)
		}
	}

	roles, err := s.Repository.GetRolesByIDs(tx, eli.CompanyID, roleIDs)
	if err != nil {
		log.WithFields(log.Fields{
			"employeeListInput": fmt.Sprintf("%+v", eli),
			"roleIDs":           roleIDs,
		}).Errorln("[Service][GetEmployeeList][GetRolesByIDs]: ", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	employeesOutput := make([]EmployeeOutput, 0)
	for _, v := range employees {
		employeesOutput = append(employeesOutput, EmployeeOutput{
			CompanyID:   v.CompanyID,
			ID:          v.ID,
			FullName:    v.FullName,
			Email:       v.Email,
			PhoneNumber: v.PhoneNumber,
			RoleID:      v.RoleID,
			RoleName:    roles[v.RoleID].Name,
			Status:      v.Status.String(),
		})
	}

	return &EmployeeListOutput{
		Count:     count,
		Employees: employeesOutput,
	}, http.StatusOK, nil
}
