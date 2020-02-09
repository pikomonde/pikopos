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

	out, status, err := h.Service.EmployeeList(in)
	if err != nil {
		log.WithFields(log.Fields{
			"in": fmt.Sprintf("%+v", in),
		}).Errorln("[Delivery][HandleEmployeeList][Login]: ", err.Error())
		respErrorJSON(w, r, status, err.Error())
		return
	}

	respSuccessJSON(w, r, status, out)
	return
}
