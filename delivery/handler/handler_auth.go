package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pikomonde/pikopos/service"
	log "github.com/sirupsen/logrus"
)

// RegisterAuth is used to register auth related handlers
func (h *Handler) RegisterAuth() {
	h.Mux.HandleFunc("/ping", ctxGET(h.HandlePing))

	h.Mux.HandleFunc("/auth/register", ctxPOST(h.HandleAuthRegister))
	h.Mux.HandleFunc("/auth/verify", ctxPOST(h.HandleAuthVerify))

	h.Mux.HandleFunc("/auth/login", ctxPOST(h.HandleAuthLogin))
	h.Mux.HandleFunc("/auth/me", ctxGET(middleAuth(h.HandleAuthMe)))
	h.Mux.HandleFunc("/auth/logout", ctxPOST(h.HandleAuthLogout))

	// h.Mux.GET("/login", h.HandlerLoginHTML)
	// h.Mux.POST("/register", h.HandlerRegisterHTML)

	// gr := h.Mux.Group("restricted")
	// gr.Use(middleAuth)
	// gr.GET("/auth/restricted", h.HandlerRestricted)
	// h.Mux.GET("/restricted", h.HandlerRestricted)
}

// HandleAuthRegister is used for user to register. User sent the register info
// (from frontend), then this API returns executes verification code
func (h *Handler) HandleAuthRegister(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	in := service.RegisterInput{}
	err := decoder.Decode(&in)
	if err != nil {
		log.WithFields(log.Fields{
			"in": fmt.Sprintf("%+v", in),
		}).Errorln("[Delivery][HandleAuthRegister][Decode]: ", err.Error())
		respErrorJSON(w, r, http.StatusBadRequest, errorWrongJSONFormat)
		return
	}

	status, err := h.Service.Register(in)
	if err != nil {
		if status == http.StatusInternalServerError {
			log.WithFields(log.Fields{
				"in": fmt.Sprintf("%+v", in),
			}).Errorln("[Delivery][HandleAuthRegister][Register]: ", err.Error())
		}
		respErrorJSON(w, r, status, errorFailedToRegister)
		return
	}

	respSuccessJSON(w, r, status, simpleMessage{"succes"})
	return
}

// HandleAuthVerify is used for user to verify their account. User sent the otp
// (from frontend), then the app change the user status to active
func (h *Handler) HandleAuthVerify(w http.ResponseWriter, r *http.Request) {
	// decoder := json.NewDecoder(r.Body)
	// in := service.LoginInput{}
	// err := decoder.Decode(&in)
	// if err != nil {
	// 	log.WithFields(log.Fields{
	// 		"in": fmt.Sprintf("%+v", in),
	// 	}).Errorln("[Delivery][HandleAuthVerify][Decode]: ", err.Error())
	// 	respErrorJSON(w, r, http.StatusBadRequest, errorWrongJSONFormat)
	// 	return
	// }

	// tokenEncoded, status, err := d.Service.Verify(li)
	// if err != nil {
	// 	ctx.JSON(status, responseAPI{
	// 		Status: status,
	// 		// ProcessTime: 0,
	// 		Data: err.Error(),
	// 	})
	// }

	// ctx.JSON(status, responseAPI{
	// 	Status: status,
	// 	// ProcessTime: 0,
	// 	Data: struct {
	// 		Token string `json:"token"`
	// 	}{tokenEncoded},
	// })
	return
}

// HandleAuthLogin is used for user to login. User sent the credentials (from
// frontend), then this API returns token
func (h *Handler) HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
	in := service.LoginInput{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&in)
	if err != nil {
		log.WithFields(log.Fields{
			"in": fmt.Sprintf("%+v", in),
		}).Errorln("[Delivery][HandleAuthLogin][Decode]: ", err.Error())
		respErrorJSON(w, r, http.StatusBadRequest, errorWrongJSONFormat)
		return
	}

	out, status, err := h.Service.Login(in)
	if err != nil {
		if status == http.StatusInternalServerError {
			log.WithFields(log.Fields{
				"in": fmt.Sprintf("%+v", in),
			}).Errorln("[Delivery][HandleAuthLogin][Login]: ", err.Error())
		}
		respErrorJSON(w, r, status, errorFailedToLogin)
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
	mid, ok := r.Context().Value(ctxKey("ctxData")).(middlewareData)
	if !ok {
		respErrorJSON(w, r, http.StatusUnauthorized, errorMissingAuthSessionData)
		return
	}

	respSuccessJSON(w, r, http.StatusOK, mid.User)
	return
}

// HandleAuthLogout is used to log out the user. It might be used if the app using
// a cache (ex: redis), it will delete the record in the cache.
func (h *Handler) HandleAuthLogout(w http.ResponseWriter, r *http.Request) {
	// mid, ok := r.Context().Value(ctxKey("ctxData")).(middlewareData)
	// if !ok {
	// 	respErrorJSON(w, r, http.StatusUnauthorized, errorMissingAuthSessionData)
	// 	return
	// }

	// respSuccessJSON(w, r, http.StatusOK, mid.User)
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
