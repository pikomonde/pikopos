package handler

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type responseAPI struct {
	Status      int         `json:"status"`
	ProcessTime int         `json:"process_time"`
	Data        interface{} `json:"data"`
}

func respErrorJSON(w http.ResponseWriter, r *http.Request, status int, errStr string) {
	processTimeRaw := r.Context().Value(ctxKey("processTime"))
	if processTimeRaw == nil {
		log.Errorln("[respSuccessJSON] processTime:",
			"don't forget to add the middleware that add 'processTime' contenxt")
		respErrorText(w, r)
		return
	}
	processTime, ok := processTimeRaw.(middlewareProcessTime)
	if !ok {
		respErrorJSON(w, r, http.StatusInternalServerError, errorMissingProcessingTimeData)
		return
	}
	js, err := json.Marshal(responseAPI{
		Status:      status,
		ProcessTime: int(processTime.ProcessTime.Sub(time.Now()) / time.Microsecond),
		Data: struct {
			Message string `json:"message"`
		}{errStr},
	})
	if err != nil {
		log.Errorln("[respSuccessJSON] marshal:", err.Error())
		respErrorText(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return
}

func respSuccessJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	processTimeRaw := r.Context().Value(ctxKey("processTime"))
	if processTimeRaw == nil {
		log.Errorln("[respSuccessJSON] processTime:",
			"don't forget to add the middleware that add 'processTime' contenxt")
		respErrorText(w, r)
		return
	}
	processTime, ok := processTimeRaw.(middlewareProcessTime)
	if !ok {
		respErrorJSON(w, r, http.StatusInternalServerError, errorMissingProcessingTimeData)
		return
	}
	js, err := json.Marshal(responseAPI{
		Status:      status,
		ProcessTime: int(time.Now().Sub(processTime.ProcessTime) / time.Microsecond),
		Data:        data,
	})
	if err != nil {
		log.Errorln("[respSuccessJSON] marshal:", err.Error())
		respErrorText(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return
}

func respErrorText(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 - Something bad happened!"))
}

type ctxKey string

const errorWrongJSONFormat = "Wrong JSON Format"

const errorWrongJWTSigningMethod = "Wrong JWT Signing Method"
const errorExpiredJWTToken = "Expired JWT Token"
const errorMissingJWTData = "Missing JWT Data"
const errorDeformedJWTToken = "Deformed JWT Token"
const errorMissingAuthSessionData = "Missing Auth Session Data"
const errorMissingProcessingTimeData = "Missing Processing Time Data"

const errorInvalidRequestMethod = "Invalid Request Method"

// config.JWTSecret
