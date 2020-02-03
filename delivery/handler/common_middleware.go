package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pikomonde/pikopos/config"
	"github.com/pikomonde/pikopos/entity"
	log "github.com/sirupsen/logrus"
	jwt "github.com/zt9/am2.piko/github.com/dgrijalva/jwt-go"
)

type middlewareProcessTime struct {
	ProcessTime time.Time
}

func ctxGET(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			respErrorJSON(w, r, http.StatusUnauthorized, errorInvalidRequestMethod)
			return
		}

		newCtx := context.WithValue(r.Context(), ctxKey("processTime"),
			middlewareProcessTime{
				ProcessTime: time.Now(),
			})
		newReq := r.WithContext(newCtx)

		next.ServeHTTP(w, newReq)
	})
}

func ctxPOST(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			respErrorJSON(w, r, http.StatusUnauthorized, errorInvalidRequestMethod)
			return
		}

		newCtx := context.WithValue(r.Context(), ctxKey("processTime"),
			middlewareProcessTime{
				ProcessTime: time.Now(),
			})
		newReq := r.WithContext(newCtx)

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
			return []byte(config.JWTSecret), nil
		})
		if err != nil {
			log.WithFields(log.Fields{
				"reqToken": reqToken,
				"token":    token,
			}).Errorln("[Delivery][middleAuth][Parse]: ", err.Error())
			respErrorJSON(w, r, http.StatusUnauthorized, errorWrongJWTSigningMethod)
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
