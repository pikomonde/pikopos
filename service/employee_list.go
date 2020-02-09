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

// EmployeeList is used to get employee list a new user
func (s *Service) EmployeeList(eli EmployeeListInput) (*EmployeeListOutput, int, error) {
	// TODO: validate input
	// TODO: change to informative error in user

	count, err := s.Repository.GetEmployeesCount(eli.CompanyID)
	if err != nil {
		log.WithFields(log.Fields{
			"registerInput": fmt.Sprintf("%+v", eli),
		}).Errorln("[Service][EmployeeList][GetEmployees]: ", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	employees, err := s.Repository.GetEmployees(eli.CompanyID, repository.Pagination{
		LastID: eli.LastID,
		Limit:  eli.Limit,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"registerInput": fmt.Sprintf("%+v", eli),
		}).Errorln("[Service][EmployeeList][GetEmployees]: ", err.Error())
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

	roles, err := s.Repository.GetRolesByIDs(eli.CompanyID, roleIDs)
	if err != nil {
		log.WithFields(log.Fields{
			"registerInput": fmt.Sprintf("%+v", eli),
			"roleIDs":       roleIDs,
		}).Errorln("[Service][EmployeeList][GetRolesByIDs]: ", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	employeesOutput := make([]EmployeeOutput, 0)
	for _, v := range employees {
		employeesOutput = append(employeesOutput, EmployeeOutput{
			CompanyID:   v.CompanyID,
			ID:          v.ID,
			FirstName:   v.FirstName,
			LastName:    v.LastName,
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
