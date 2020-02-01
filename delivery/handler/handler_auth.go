package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pikomonde/pikopos/config"
	"github.com/pikomonde/pikopos/entity"
	"github.com/pikomonde/pikopos/service"
	log "github.com/sirupsen/logrus"
)

// Handler is used
type Handler struct {
	Service *service.Service
	Mux     *http.ServeMux
}

type middlewareData struct {
	User entity.ServiceUserSession
}

// RegisterAuth is used to register auth related handlers
func (h *Handler) RegisterAuth() {
	h.Mux.HandleFunc("/ping", h.HandlePing)

	// h.Mux.HandleFunc("/auth/register", h.HandleAuthRegister)
	// h.Mux.HandleFunc("/auth/verify", h.HandleAuthVerify)

	h.Mux.HandleFunc("/auth/login", h.HandleAuthLogin)
	h.Mux.HandleFunc("/auth/me", middleAuth(h.HandleAuthMe))

	// h.Mux.GET("/login", h.HandlerLoginHTML)
	// h.Mux.POST("/register", h.HandlerRegisterHTML)

	// gr := h.Mux.Group("restricted")
	// gr.Use(middleAuth)
	// gr.GET("/auth/restricted", h.HandlerRestricted)
	// h.Mux.GET("/restricted", h.HandlerRestricted)
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

		newCtx := context.WithValue(nil, ctxData("ctxData"), middlewareData{
			User: userSession,
		})
		newReq := r.WithContext(newCtx)

		next.ServeHTTP(w, newReq)
	})
}

// HandleAuthRegister is used for user to register. User sent the register info
// (from frontend), then this API returns executes verification code
func (h *Handler) HandleAuthRegister(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	in := service.RegisterInput{}
	err := decoder.Decode(&in)
	if err != nil {
		log.WithFields(log.Fields{}).
			Errorln("[Delivery][HandleAuthRegister][Decode]: ", err.Error())
		respErrorJSON(w, r, http.StatusBadRequest, errorWrongJSONFormat)
		return
	}

	status, err := h.Service.Register(in)
	if err != nil {
		log.WithFields(log.Fields{
			"in": fmt.Sprintf("%+v", in),
		}).Errorln("[Delivery][HandleAuthRegister][Login]: ", err.Error())
		respErrorJSON(w, r, status, err.Error())
		return
	}

	respSuccessJSON(w, r, status, nil)
	return
}

// func (h *Handler) HandleAuthVerify(w http.ResponseWriter, req *http.Request) {
// 	decoder := json.NewDecoder(ctx.Request().Body)
// 	li := service.LoginInput{}
// 	err := decoder.Decode(&li)
// 	if err != nil {
// 		log.WithFields(log.Fields{}).Errorln("[Delivery][HandleAuthVerify][Decode]: ", err.Error())
// 		ctx.JSON(http.StatusBadRequest, responseAPI{
// 			Status: http.StatusBadRequest,
// 			// ProcessTime: 0,
// 			Data: errors.New(errorWrongJSONFormat).Error(),
// 		})
// 		return err
// 	}

// 	tokenEncoded, status, err := d.Service.Login(li)
// 	if err != nil {
// 		ctx.JSON(status, responseAPI{
// 			Status: status,
// 			// ProcessTime: 0,
// 			Data: err.Error(),
// 		})
// 	}

// 	// cookie := new(http.Cookie)
// 	// cookie.Name = "authentication"
// 	// cookie.Value = tokenEncoded
// 	// cookie.Expires = time.Now().Add(72 * time.Hour)
// 	// cookie.Domain = "localhost"
// 	// cookie.Path = "/"
// 	// // TODO: enable secure for production
// 	// // cookie.Secure = true
// 	// ctx.SetCookie(cookie)

// 	ctx.JSON(status, responseAPI{
// 		Status: status,
// 		// ProcessTime: 0,
// 		Data: struct {
// 			Token string `json:"token"`
// 		}{tokenEncoded},
// 	})
// 	return nil
// }

// HandleAuthLogin is used for user to login. User sent the credentials (from
// frontend), then this API returns token
func (h *Handler) HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
	in := service.LoginInput{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&in)
	if err != nil {
		log.WithFields(log.Fields{}).
			Errorln("[Delivery][HandleAuthLogin][Decode]: ", err.Error())
		respErrorJSON(w, r, http.StatusBadRequest, errorWrongJSONFormat)
		return
	}

	out, status, err := h.Service.Login(in)
	if err != nil {
		log.WithFields(log.Fields{
			"in": fmt.Sprintf("%+v", in),
		}).Errorln("[Delivery][HandleAuthLogin][Login]: ", err.Error())
		respErrorJSON(w, r, status, err.Error())
		return
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

	respSuccessJSON(w, r, status, out)
	return
}

// HandleAuthMe is used to get user information from token. It should use only
// to verify and extract token data. This endpoint can be also done in the
// frontend, but the data will not be verified. It is in the frontend anyway.
func (h *Handler) HandleAuthMe(w http.ResponseWriter, r *http.Request) {
	mid, ok := r.Context().Value(ctxData("ctxData")).(middlewareData)
	if !ok {
		respErrorJSON(w, r, http.StatusUnauthorized, errorMissingAuthSessionData)
		return
	}

	respSuccessJSON(w, r, http.StatusOK, mid.User)
	return
}

// HandlePing is used for pinging, mainly for test whether the server alive or
// not
func (h *Handler) HandlePing(w http.ResponseWriter, r *http.Request) {
	respSuccessJSON(w, r, http.StatusOK, struct {
		Message string `json:"message"`
	}{"pong"})
	return
}

// func (h *Handler) HandlerRestricted(w http.ResponseWriter, req *http.Request) error {
// 	// ctx.Get()

// 	ctx.HTML(200, fmt.Sprintf("-----"))
// 	return nil
// }

// func (h *Handler) HandlerLoginHTML(w http.ResponseWriter, req *http.Request) error {
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
