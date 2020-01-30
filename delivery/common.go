package delivery

import (
	"fmt"

	"github.com/labstack/echo"
)

type responseAPI struct {
	Status      int         `json:"status"`
	ProcessTime int         `json:"process_time"`
	Data        interface{} `json:"data"`
}

func respErrorJSON(ctx echo.Context, status int, errStr string) error {
	ctx.JSON(status, responseAPI{
		Status: status,
		// ProcessTime: 0,
		Data: errStr,
	})
	return fmt.Errorf(errStr)
}

func respSuccessJSON(ctx echo.Context, status int, data interface{}) error {
	ctx.JSON(status, responseAPI{
		Status: status,
		// ProcessTime: 0,
		Data: data,
	})
	return nil
}

const errorWrongJSONFormat = "Wrong JSON Format"

const errorWrongJWTSigningMethod = "Wrong JWT Signing Method"
const errorExpiredJWTToken = "Expired JWT Token"
const errorMissingJWTData = "Missing JWT Data"
const errorDeformedJWTToken = "Deformed JWT Token"
const errorMissingAuthSessionData = "Missing Auth Session Data"

// config.JWTSecret
