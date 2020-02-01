package handler

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type responseAPI struct {
	Status      int         `json:"status"`
	ProcessTime int         `json:"process_time"`
	Data        interface{} `json:"data"`
}

func respErrorJSON(w http.ResponseWriter, r *http.Request, status int, errStr string) {
	// TODO: do ProcessTime
	js, err := json.Marshal(responseAPI{
		Status: status,
		// ProcessTime: 0,
		Data: struct {
			Message string `json:"message"`
		}{errStr},
	})
	if err != nil {
		w.Header().Set("Content-Type", "text")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		log.WithFields(log.Fields{}).Errorln("[respSuccessJSON] marshal:", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return
}

func respSuccessJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	// TODO: do ProcessTime
	js, err := json.Marshal(responseAPI{
		Status: status,
		// ProcessTime: 0,
		Data: data,
	})
	if err != nil {
		w.Header().Set("Content-Type", "text")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		log.WithFields(log.Fields{}).Errorln("[respSuccessJSON] marshal:", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return
}

type ctxData string

const errorWrongJSONFormat = "Wrong JSON Format"

const errorWrongJWTSigningMethod = "Wrong JWT Signing Method"
const errorExpiredJWTToken = "Expired JWT Token"
const errorMissingJWTData = "Missing JWT Data"
const errorDeformedJWTToken = "Deformed JWT Token"
const errorMissingAuthSessionData = "Missing Auth Session Data"

// config.JWTSecret
