package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/pikomonde/pikopos/config"
	"github.com/pikomonde/pikopos/entity"
	"github.com/pikomonde/pikopos/service"
	log "github.com/sirupsen/logrus"
)

type middlewareData struct {
	User entity.ServiceUserSession
}

// PrepareAuth is used to prepare auth related endpoints
func (d *Delivery) PrepareAuth() {
	d.EchoDelivery.POST("/auth/register", d.HandleAuthRegister)
	d.EchoDelivery.POST("/auth/verify", d.HandleAuthVerify)

	d.EchoDelivery.POST("/auth/login", d.HandleAuthLogin)
	d.EchoDelivery.GET("/auth/me", middleAuth(d.HandleAuthMe))

	// d.EchoDelivery.GET("/login", d.HandlerLoginHTML)
	// d.EchoDelivery.POST("/register", d.HandlerRegisterHTML)

	// gr := d.EchoDelivery.Group("restricted")
	// gr.Use(middleAuth)
	// gr.GET("/auth/restricted", d.HandlerRestricted)
	// d.EchoDelivery.GET("/restricted", d.HandlerRestricted)
}

func middleAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		reqToken := ctx.Request().Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		if len(splitToken) != 2 {
			return respErrorJSON(ctx, http.StatusUnauthorized, errorDeformedJWTToken)
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
			return respErrorJSON(ctx, http.StatusUnauthorized, errorWrongJWTSigningMethod)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return respErrorJSON(ctx, http.StatusUnauthorized, errorExpiredJWTToken)
		}

		userSession := entity.ServiceUserSession{}
		userSessionStr, ok := claims["data"].(string)
		if !ok {
			return respErrorJSON(ctx, http.StatusUnauthorized, errorMissingJWTData)
		}

		err = json.Unmarshal([]byte(userSessionStr), &userSession)
		if err != nil {
			return respErrorJSON(ctx, http.StatusUnauthorized, errorDeformedJWTToken)
		}

		ctx.Set("data", middlewareData{
			User: userSession,
		})
		return next(ctx)
	}
}

func (d *Delivery) HandleAuthRegister(ctx echo.Context) error {
	// TODO: process time
	decoder := json.NewDecoder(ctx.Request().Body)
	ri := service.RegisterInput{}
	err := decoder.Decode(&ri)
	if err != nil {
		log.WithFields(log.Fields{}).Errorln("[Delivery][HandleAuthRegister][Decode]: ", err.Error())
		ctx.JSON(http.StatusBadRequest, responseAPI{
			Status: http.StatusBadRequest,
			// ProcessTime: 0,
			Data: errors.New(errorWrongJSONFormat).Error(),
		})
		return err
	}

	status, err := d.Service.Register(ri)
	if err != nil {
		ctx.JSON(status, responseAPI{
			Status: status,
			// ProcessTime: 0,
			Data: err.Error(),
		})
	}
	ctx.JSON(status, responseAPI{
		Status: status,
		// ProcessTime: 0,
		Data: nil,
	})
	return nil
}

func (d *Delivery) HandleAuthVerify(ctx echo.Context) error {
	decoder := json.NewDecoder(ctx.Request().Body)
	li := service.LoginInput{}
	err := decoder.Decode(&li)
	if err != nil {
		log.WithFields(log.Fields{}).Errorln("[Delivery][HandleAuthVerify][Decode]: ", err.Error())
		ctx.JSON(http.StatusBadRequest, responseAPI{
			Status: http.StatusBadRequest,
			// ProcessTime: 0,
			Data: errors.New(errorWrongJSONFormat).Error(),
		})
		return err
	}

	tokenEncoded, status, err := d.Service.Login(li)
	if err != nil {
		ctx.JSON(status, responseAPI{
			Status: status,
			// ProcessTime: 0,
			Data: err.Error(),
		})
	}

	// cookie := new(http.Cookie)
	// cookie.Name = "authentication"
	// cookie.Value = tokenEncoded
	// cookie.Expires = time.Now().Add(72 * time.Hour)
	// cookie.Domain = "localhost"
	// cookie.Path = "/"
	// // TODO: enable secure for production
	// // cookie.Secure = true
	// ctx.SetCookie(cookie)

	ctx.JSON(status, responseAPI{
		Status: status,
		// ProcessTime: 0,
		Data: struct {
			Token string `json:"token"`
		}{tokenEncoded},
	})
	return nil
}

func (d *Delivery) HandleAuthLogin(ctx echo.Context) error {
	decoder := json.NewDecoder(ctx.Request().Body)
	li := service.LoginInput{}
	err := decoder.Decode(&li)
	if err != nil {
		log.WithFields(log.Fields{}).Errorln("[Delivery][HandleAuthLogin][Decode]: ", err.Error())
		return respErrorJSON(ctx, http.StatusBadRequest, errors.New(errorWrongJSONFormat).Error())
	}

	tokenEncoded, status, err := d.Service.Login(li)
	if err != nil {
		// TODO log error
		return respErrorJSON(ctx, status, err.Error())
	}

	// cookie := new(http.Cookie)
	// cookie.Name = "authentication"
	// cookie.Value = tokenEncoded
	// cookie.Expires = time.Now().Add(72 * time.Hour)
	// cookie.Domain = "localhost"
	// cookie.Path = "/"
	// // TODO: enable secure for production
	// // cookie.Secure = true
	// ctx.SetCookie(cookie)

	return respSuccessJSON(ctx, status, struct {
		Token string `json:"token"`
	}{tokenEncoded})
}

func (d *Delivery) HandleAuthMe(ctx echo.Context) error {
	mid, ok := ctx.Get("data").(middlewareData)
	if !ok {
		return respErrorJSON(ctx, http.StatusUnauthorized, errorMissingAuthSessionData)
	}

	return respSuccessJSON(ctx, http.StatusOK, mid.User)
}

// func (d *Delivery) HandlerRestricted(ctx echo.Context) error {
// 	// ctx.Get()

// 	ctx.HTML(200, fmt.Sprintf("-----"))
// 	return nil
// }

// func (d *Delivery) HandlerLoginHTML(ctx echo.Context) error {
// 	// ctx.Get()

// 	ctx.HTML(200, `
// 	<script>
// 	  function submit(){
// 			var xhr = new XMLHttpRequest();
// 			var url = "http://localhost:1235/auth/login";
// 			xhr.open("POST", url, true);
// 			xhr.setRequestHeader("Content-Type", "application/json");
// 			xhr.onreadystatechange = function () {
// 					if (xhr.readyState === 4 && xhr.status === 200) {
// 							var json = JSON.parse(xhr.responseText);
// 							console.log(json.email + ", " + json.password);
// 					}
// 			};
// 			var data = JSON.stringify({"company_username": "toko_bunga_03", "employee_identifier": "nugraha@gmail.com", "password": "ABC5dasar"});
// 			xhr.send(data);

// 		}
// 	</script>
// 	<input type="text" id="user"/><br/>
// 	<input type="password" id="pass"/><br/>
// 	<button onclick="submit()">Hohoho</button>
// 	`)
// 	return nil
// }
