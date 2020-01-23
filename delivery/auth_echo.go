package delivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/pikomonde/pikopos/service"
	log "github.com/sirupsen/logrus"
)

// PrepareAuth is used to prepare auth related endpoints
func (d *Delivery) PrepareAuth() {
	d.EchoDelivery.POST("/auth/register", d.HandlerRegister)
	d.EchoDelivery.POST("/auth/login", d.HandlerLogin)
	d.EchoDelivery.POST("/auth/verify", d.HandlerVerify)

	// d.EchoDelivery.GET("/login", d.HandlerLoginHTML)
	// d.EchoDelivery.POST("/register", d.HandlerRegisterHTML)

	gr := d.EchoDelivery.Group("restricted")
	gr.Use(middleAuth)
	// gr.GET("/auth/restricted", d.HandlerRestricted)
	// d.EchoDelivery.GET("/restricted", d.HandlerRestricted)
}

func middleAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		_, err := ctx.Cookie("authentication")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, responseAPI{
				Status: http.StatusUnauthorized,
				// ProcessTime: 0,
				Data: nil,
			})
			return err
		}
		return next(ctx)
	}
}

func (d *Delivery) HandlerRegister(ctx echo.Context) error {
	// TODO: process time
	decoder := json.NewDecoder(ctx.Request().Body)
	ri := service.RegisterInput{}
	err := decoder.Decode(&ri)
	if err != nil {
		log.WithFields(log.Fields{}).Errorln("[Delivery][HandlerRegister][Decode]: ", err.Error())
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

func (d *Delivery) HandlerLogin(ctx echo.Context) error {
	decoder := json.NewDecoder(ctx.Request().Body)
	li := service.LoginInput{}
	err := decoder.Decode(&li)
	if err != nil {
		log.WithFields(log.Fields{}).Errorln("[Delivery][HandlerLogin][Decode]: ", err.Error())
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

	cookie := new(http.Cookie)
	cookie.Name = "authentication"
	cookie.Value = tokenEncoded
	cookie.Expires = time.Now().Add(72 * time.Hour)
	cookie.Domain = "localhost"
	cookie.Path = "/"
	// TODO: enable secure for production
	// cookie.Secure = true
	ctx.SetCookie(cookie)

	ctx.JSON(status, responseAPI{
		Status: status,
		// ProcessTime: 0,
		Data: nil,
	})
	return nil
}

func (d *Delivery) HandlerVerify(ctx echo.Context) error {
	decoder := json.NewDecoder(ctx.Request().Body)
	li := service.LoginInput{}
	err := decoder.Decode(&li)
	if err != nil {
		log.WithFields(log.Fields{}).Errorln("[Delivery][HandlerLogin][Decode]: ", err.Error())
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

	cookie := new(http.Cookie)
	cookie.Name = "authentication"
	cookie.Value = tokenEncoded
	cookie.Expires = time.Now().Add(72 * time.Hour)
	cookie.Domain = "localhost"
	cookie.Path = "/"
	// TODO: enable secure for production
	// cookie.Secure = true
	ctx.SetCookie(cookie)

	ctx.JSON(status, responseAPI{
		Status: status,
		// ProcessTime: 0,
		Data: nil,
	})
	return nil
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
