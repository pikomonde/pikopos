package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pikomonde/pikopos/service"
	log "github.com/sirupsen/logrus"
)

// RegisterEmployee is used to register employee related handlers
func (h *Handler) RegisterEmployee() {
	h.Mux.HandleFunc("/employee/list", ctxGET(middleAuth(h.HandleEmployeeList)))
	h.Mux.HandleFunc("/employee/create", ctxPOST(middleAuth(h.HandleEmployeeCreate)))
	h.Mux.HandleFunc("/employee/update", ctxPOST(middleAuth(h.HandleEmployeeUpdate)))
}

// HandleEmployeeList is used to list all employee in the same company
func (h *Handler) HandleEmployeeList(w http.ResponseWriter, r *http.Request) {
	mid, ok := r.Context().Value(ctxKey("ctxData")).(middlewareData)
	if !ok {
		respErrorJSON(w, r, http.StatusUnauthorized, errorMissingAuthSessionData)
		return
	}

	decoder := json.NewDecoder(r.Body)
	in := service.EmployeeListInput{}
	err := decoder.Decode(&in)
	if err != nil {
		log.WithFields(log.Fields{
			"in": fmt.Sprintf("%+v", in),
		}).Errorln("[Delivery][HandleEmployeeList][Decode]: ", err.Error())
		respErrorJSON(w, r, http.StatusBadRequest, errorWrongJSONFormat)
		return
	}
	in.CompanyID = mid.User.CompanyID

	out, status, err := h.Service.GetEmployeeList(in)
	if err != nil {
		if status == http.StatusInternalServerError {
			log.WithFields(log.Fields{
				"in": fmt.Sprintf("%+v", in),
			}).Errorln("[Delivery][HandleEmployeeList][GetEmployeeList]: ", err.Error())
		}
		respErrorJSON(w, r, status, err.Error())
		return
	}

	respSuccessJSON(w, r, status, out)
	return
}

// HandleEmployeeCreate is used to create a new employee in the same company
func (h *Handler) HandleEmployeeCreate(w http.ResponseWriter, r *http.Request) {
	mid, ok := r.Context().Value(ctxKey("ctxData")).(middlewareData)
	if !ok {
		respErrorJSON(w, r, http.StatusUnauthorized, errorMissingAuthSessionData)
		return
	}

	decoder := json.NewDecoder(r.Body)
	in := service.EmployeeCreateInput{}
	err := decoder.Decode(&in)
	if err != nil {
		log.WithFields(log.Fields{
			"in": fmt.Sprintf("%+v", in),
		}).Errorln("[Delivery][HandleEmployeeCreate][Decode]: ", err.Error())
		respErrorJSON(w, r, http.StatusBadRequest, errorWrongJSONFormat)
		return
	}
	in.CompanyID = mid.User.CompanyID

	out, status, err := h.Service.CreateEmployee(in)
	if err != nil {
		if status == http.StatusInternalServerError {
			log.WithFields(log.Fields{
				"in": fmt.Sprintf("%+v", in),
			}).Errorln("[Delivery][HandleEmployeeCreate][CreateEmployee]: ", err.Error())
		}
		respErrorJSON(w, r, status, err.Error())
		return
	}

	respSuccessJSON(w, r, status, out)
	return
}

// HandleEmployeeUpdate is used to update an existing employee in the same company
func (h *Handler) HandleEmployeeUpdate(w http.ResponseWriter, r *http.Request) {
	mid, ok := r.Context().Value(ctxKey("ctxData")).(middlewareData)
	if !ok {
		respErrorJSON(w, r, http.StatusUnauthorized, errorMissingAuthSessionData)
		return
	}

	decoder := json.NewDecoder(r.Body)
	in := service.EmployeeUpdateInput{}
	err := decoder.Decode(&in)
	if err != nil {
		log.WithFields(log.Fields{
			"in": fmt.Sprintf("%+v", in),
		}).Errorln("[Delivery][HandleEmployeeUpdate][Decode]: ", err.Error())
		respErrorJSON(w, r, http.StatusBadRequest, errorWrongJSONFormat)
		return
	}
	in.CompanyID = mid.User.CompanyID

	out, status, err := h.Service.UpdateEmployee(in)
	if err != nil {
		if status == http.StatusInternalServerError {
			log.WithFields(log.Fields{
				"in": fmt.Sprintf("%+v", in),
			}).Errorln("[Delivery][HandleEmployeeUpdate][UpdateEmployee]: ", err.Error())
		}
		respErrorJSON(w, r, status, err.Error())
		return
	}

	respSuccessJSON(w, r, status, out)
	return
}
