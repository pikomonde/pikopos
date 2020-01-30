package service

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

type LoginInput struct {
	CompanyUsername    string `max:"16" json:"company_username"`
	EmployeeIdentifier string `max:"48" json:"employee_identifier"`
	Password           string `min:"8" json:"password"`
}

func (s *Service) Login(li LoginInput) (string, int, error) {
	// TODO: validate input
	// TODO: change to informative error in user

	// remove password for security while logging
	passwordRaw := li.Password
	li.Password = ""

	company, err := s.Repository.GetCompanyByUsername(li.CompanyUsername)
	if err != nil {
		log.WithFields(log.Fields{
			"loginInput": fmt.Sprintf("%+v", li),
		}).Errorln("[Service][Login][GetCompanyByUsername]: ", err.Error())
		return "", http.StatusInternalServerError, err
	}

	employee, err := s.Repository.GetEmployeeByIdentifier(company.ID, li.EmployeeIdentifier)
	if err != nil {
		log.WithFields(log.Fields{
			"loginInput": fmt.Sprintf("%+v", li),
			"company":    fmt.Sprintf("%+v", company),
		}).Errorln("[Service][Login][GetEmployeeByIdentifier]: ", err.Error())
		return "", http.StatusInternalServerError, err
	}

	passwordHashed := common.SHA256(fmt.Sprintf("%s-%s-%s-%d",
		passwordRaw, employee.Email, employee.PhoneNumber, employee.ID))
	expectedPasswordHashed, err := s.Repository.GetEmployeePassword(employee.CompanyID, employee.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"loginInput": fmt.Sprintf("%+v", li),
			"company":    fmt.Sprintf("%+v", company),
			"employee":   fmt.Sprintf("%+v", employee),
		}).Errorln("[Service][Login][GetEmployeePassword]: ", err.Error())
		return "", http.StatusInternalServerError, err
	}
	if passwordHashed != expectedPasswordHashed {
		return "", http.StatusUnauthorized, errors.New("Wrong Email/Phonenumber or Password")
	}

	// TODO: get privileges
	// role, err := s.Repository.CreateRole(entity.Role{
	// 	CompanyID: company.ID,
	// 	Name:      entity.RoleSuperAdmin,
	// 	Status:    entity.RoleStatusActive,
	// })
	// if err != nil {
	// 	log.WithFields(log.Fields{
	// 		"registerInput": fmt.Sprintf("%+v", ri),
	// 		"company":       fmt.Sprintf("%+v", company),
	// 	}).Errorln("[Service][Register][CreateRole]: ", err.Error())
	// 	return http.StatusInternalServerError, err
	// }

	serviceUserSession := entity.ServiceUserSession{
		CompanyID:       company.ID,
		CompanyUsername: company.Username,
		CompanyName:     company.Name,
		ID:              employee.ID,
		FullName:        fmt.Sprintf("%s %s", employee.FirstName, employee.LastName),
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
		}).Errorln("[Service][Login][Marshal serviceUserSession]: ", err.Error())
		return "", http.StatusInternalServerError, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": string(serviceUserSessionBytes),
		"exp":  time.Now().Add(tokenCookieExpire).Unix(),
	})
	tokenEncoded, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		log.WithFields(log.Fields{
			"loginInput":         fmt.Sprintf("%+v", li),
			"company":            fmt.Sprintf("%+v", company),
			"employee":           fmt.Sprintf("%+v", employee),
			"serviceUserSession": fmt.Sprintf("%+v", serviceUserSession),
		}).Errorln("[Service][Login][SignedString]: ", err.Error())
		return "", http.StatusInternalServerError, err
	}

	return tokenEncoded, http.StatusOK, nil
}
