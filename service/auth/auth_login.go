package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pikomonde/pikopos/common"
	"github.com/pikomonde/pikopos/config"
	"github.com/pikomonde/pikopos/entity"
	log "github.com/sirupsen/logrus"
)

const (
	tokenCookieExpire = 72 * time.Hour
)

// LoginInput is used as request for login
type LoginInput struct {
	CompanyUsername    string `max:"16" json:"company_username"`
	EmployeeIdentifier string `max:"48" json:"employee_identifier"`
	Password           string `min:"8" json:"password"`
}

// LoginOutput is used as response for login
type LoginOutput struct {
	Token string `json:"token"`
}

// Login is used to logged in a user
func (s *ServiceAuth) Login(li LoginInput) (*LoginOutput, int, error) {
	// TODO: validate input
	// TODO: change to informative error in user

	// remove password for security while logging
	passwordRaw := li.Password
	li.Password = ""

	company, err := s.RepositoryCompany.GetCompanyByUsername(nil, li.CompanyUsername)
	if err != nil {
		log.WithFields(log.Fields{
			"loginInput": fmt.Sprintf("%+v", li),
		}).Errorln("[ServiceAuth][Login][GetCompanyByUsername]: ", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	employee, err := s.RepositoryEmployee.GetEmployeeByIdentifier(nil, company.ID, li.EmployeeIdentifier)
	if err != nil {
		log.WithFields(log.Fields{
			"loginInput": fmt.Sprintf("%+v", li),
			"company":    fmt.Sprintf("%+v", company),
		}).Errorln("[ServiceAuth][Login][GetEmployeeByIdentifier]: ", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	passwordHashed := common.SHA256(fmt.Sprintf("%s-%s-%s-%d",
		passwordRaw, employee.Email, employee.PhoneNumber, employee.ID))
	expectedPasswordHashed, err := s.RepositoryEmployee.GetEmployeePassword(nil, employee.CompanyID, employee.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"loginInput": fmt.Sprintf("%+v", li),
			"company":    fmt.Sprintf("%+v", company),
			"employee":   fmt.Sprintf("%+v", employee),
		}).Errorln("[ServiceAuth][Login][GetEmployeePassword]: ", err.Error())
		return nil, http.StatusInternalServerError, err
	}
	if passwordHashed != expectedPasswordHashed {
		return nil, http.StatusUnauthorized, errors.New(common.ErrorWrongLoginInfo)
	}

	// TODO: get privileges
	// role, err := s.RepositoryRole.CreateRole(nil, entity.Role{
	// 	CompanyID: company.ID,
	// 	Name:      entity.RoleSuperAdmin,
	// 	Status:    entity.RoleStatusActive,
	// })
	// if err != nil {
	// 	log.WithFields(log.Fields{
	// 		"registerInput": fmt.Sprintf("%+v", ri),
	// 		"company":       fmt.Sprintf("%+v", company),
	// 	}).Errorln("[ServiceAuth][Register][CreateRole]: ", err.Error())
	// 	return http.StatusInternalServerError, err
	// }

	serviceUserSession := entity.ServiceUserSession{
		CompanyID:       company.ID,
		CompanyUsername: company.Username,
		CompanyName:     company.Name,
		ID:              employee.ID,
		FullName:        employee.FullName,
		Email:           employee.Email,
		PhoneNumber:     employee.PhoneNumber,
		Privileges:      []string{},
	}
	serviceUserSessionBytes, err := json.Marshal(serviceUserSession)
	if err != nil {
		log.WithFields(log.Fields{
			"loginInput":         fmt.Sprintf("%+v", li),
			"company":            fmt.Sprintf("%+v", company),
			"employee":           fmt.Sprintf("%+v", employee),
			"serviceUserSession": fmt.Sprintf("%+v", serviceUserSession),
		}).Errorln("[ServiceAuth][Login][Marshal serviceUserSession]: ", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": string(serviceUserSessionBytes),
		"exp":  time.Now().Add(tokenCookieExpire).Unix(),
	})
	tokenEncoded, err := token.SignedString([]byte(config.C.JWTSecret))
	if err != nil {
		log.WithFields(log.Fields{
			"loginInput":         fmt.Sprintf("%+v", li),
			"company":            fmt.Sprintf("%+v", company),
			"employee":           fmt.Sprintf("%+v", employee),
			"serviceUserSession": fmt.Sprintf("%+v", serviceUserSession),
		}).Errorln("[ServiceAuth][Login][SignedString]: ", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return &LoginOutput{
		tokenEncoded,
	}, http.StatusOK, nil
}
