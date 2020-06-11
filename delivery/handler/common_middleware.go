package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pikomonde/pikopos/config"
	"github.com/pikomonde/pikopos/entity"
	log "github.com/sirupsen/logrus"
)

type middlewareProcessTime struct {
	ProcessTime time.Time
}

func ctxGET(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), ctxKey("processTime"),
			middlewareProcessTime{
				ProcessTime: time.Now(),
			})
		newReq := r.WithContext(newCtx)

		// TODO: this OPTIONS is for cors only
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			return
		}
		if r.Method != "GET" {
			respErrorJSON(w, newReq, http.StatusMethodNotAllowed, errorInvalidRequestMethod)
			return
		}

		next.ServeHTTP(w, newReq)
	})
}

func ctxPOST(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), ctxKey("processTime"),
			middlewareProcessTime{
				ProcessTime: time.Now(),
			})
		newReq := r.WithContext(newCtx)

		// TODO: this OPTIONS is for cors only
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			return
		}
		if r.Method != "POST" {
			respErrorJSON(w, newReq, http.StatusMethodNotAllowed, errorInvalidRequestMethod)
			return
		}

		next.ServeHTTP(w, newReq)
	})
}

type middlewareData struct {
	User entity.ServiceUserSession
}

func middleAuth(next http.HandlerFunc) http.HandlerFunc {
	// TODO: Optimize this if possible
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		if len(splitToken) != 2 {
			respErrorJSON(w, r, http.StatusUnauthorized, errorDeformedJWTToken)
			return
		}
		reqToken = splitToken[1]

		token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
			// validation
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.C.JWTSecret), nil
		})
		if err != nil {
			log.WithFields(log.Fields{
				"reqToken": reqToken,
				"token":    token,
			}).Infoln("[Delivery][middleAuth][Parse]: ", err.Error())
			respErrorJSON(w, r, http.StatusUnauthorized, errorCredentialProblem)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			respErrorJSON(w, r, http.StatusUnauthorized, errorExpiredJWTToken)
			return
		}

		userSession := entity.ServiceUserSession{}
		userSessionStr, ok := claims["data"].(string)
		if !ok {
			respErrorJSON(w, r, http.StatusUnauthorized, errorMissingJWTData)
			return
		}

		err = json.Unmarshal([]byte(userSessionStr), &userSession)
		if err != nil {
			respErrorJSON(w, r, http.StatusUnauthorized, errorDeformedJWTToken)
			return
		}

		newCtx := context.WithValue(r.Context(), ctxKey("ctxData"),
			middlewareData{
				User: userSession,
			})
		newReq := r.WithContext(newCtx)

		next.ServeHTTP(w, newReq)
	})
}

func setProvider(next http.HandlerFunc, provider string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), ctxKey("ctxProvider"), provider)
		newReq := r.WithContext(newCtx)

		next.ServeHTTP(w, newReq)
	})
}
