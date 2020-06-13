package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/pikomonde/pikopos/config"
	log "github.com/sirupsen/logrus"
)

// RegisterAuth is used to register auth related handlers
func (h *Handler) RegisterAuth() {
	h.Mux.HandleFunc("/ping", ctxGET(h.HandlePing))

	// h.Mux.HandleFunc("/auth/login", ctxPOST(h.HandleAuthLogin))
	h.Mux.HandleFunc("/auth/me", ctxGET(middleAuth(h.HandleAuthMe)))
	h.Mux.HandleFunc("/auth/logout", ctxPOST(h.HandleAuthLogout))

	h.Mux.HandleFunc("/auth/login", ctxGET(h.HandleAuthProviderLogin))
	h.Mux.HandleFunc("/auth/callback", ctxGET(h.HandleAuthProviderCallback))

	// h.Mux.GET("/login", h.HandlerLoginHTML)
	// h.Mux.POST("/register", h.HandlerRegisterHTML)

	// gr := h.Mux.Group("restricted")
	// gr.Use(middleAuth)
	// gr.GET("/auth/restricted", h.HandlerRestricted)
	// h.Mux.GET("/restricted", h.HandlerRestricted)
}

// // HandleAuthLogin is used for user to login. User sent the credentials (from
// // frontend), then this API returns token
// func (h *Handler) HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
// 	in := service.LoginInput{}
// 	decoder := json.NewDecoder(r.Body)

// 	err := decoder.Decode(&in)
// 	if err != nil {
// 		log.WithFields(log.Fields{
// 			"in": fmt.Sprintf("%+v", in),
// 		}).Errorln("[Delivery][HandleAuthLogin][Decode]: ", err.Error())
// 		respErrorJSON(w, r, http.StatusBadRequest, errorWrongJSONFormat)
// 		return
// 	}

// 	out, status, err := h.ServiceAuth.Login(in)
// 	if err != nil {
// 		if status == http.StatusInternalServerError {
// 			log.WithFields(log.Fields{
// 				"in": fmt.Sprintf("%+v", in),
// 			}).Errorln("[Delivery][HandleAuthLogin][Login]: ", err.Error())
// 		}
// 		respErrorJSON(w, r, status, errorFailedToLogin)
// 		return
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

// 	respSuccessJSON(w, r, status, out)
// 	return
// }

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

// HandleAuthProviderLogin is used for user to login using identity provider such as
// google, facebook, twitter, etc. This API will redirect user to provider's login
// page.
func (h *Handler) HandleAuthProviderLogin(w http.ResponseWriter, r *http.Request) {
	provider := r.FormValue("provider")

	// getting auth state and redirectURL
	state, redirectURL, err := h.ServiceAuth.GenerateStateAndGetAuthURL(provider)
	if err != nil {
		log.WithFields(log.Fields{
			"provider": provider,
		}).Errorln("[Handler][HandleAuthProviderLogin][GenerateStateAndGetAuthURL]: ", err.Error())
	}

	// setting up cookie
	cookie := new(http.Cookie)
	cookie.Name = "oauthstate"
	cookie.Value = state
	cookie.Expires = time.Now().Add(10 * time.Minute)
	cookie.Domain = config.C.BaseURL
	cookie.Path = "/"
	// enable secure cookie for production
	if os.Getenv("env") == "PROD" {
		cookie.Secure = true
	}
	http.SetCookie(w, cookie)

	// redirect to provider
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)

	return
}

// HandleAuthProviderCallback is used in login proses. This endpoint specifically
// used for the provider to give state and code. This endpoint will use the state to
// validate from CSRF attack and the code is used for getting the token from
// provider through the exchange process.
func (h *Handler) HandleAuthProviderCallback(w http.ResponseWriter, r *http.Request) {
	provider := r.FormValue("provider")

	// get state from cookie
	stateCookie, err := r.Cookie("oauthstate")
	if err != nil {
		log.WithFields(log.Fields{
			"provider":    provider,
			"stateCookie": stateCookie,
		}).Errorln("[Handler][HandleAuthProviderCallback][Cookie State]: ", err.Error())
		return
	}
	stateFromCookie := stateCookie.Value

	// get state from provider
	stateFromProvider := r.FormValue("state")

	// validate state
	if stateFromCookie != stateFromProvider {
		log.WithFields(log.Fields{
			"provider":          provider,
			"stateCookie":       stateCookie,
			"stateFromCookie":   stateFromCookie,
			"stateFromProvider": stateFromProvider,
		}).Errorln("[Handler][HandleAuthProviderCallback][ValidateState]: ", err.Error())
		// TODO: redirect to front-end not authorized
		return
	}

	// getting provider client using code
	code := r.FormValue("code")

	token, err := h.ServiceAuth.Exchange(provider, code)
	if err != nil {
		log.WithFields(log.Fields{
			"provider":          provider,
			"stateCookie":       stateCookie,
			"stateFromCookie":   stateFromCookie,
			"stateFromProvider": stateFromProvider,
			"code":              code,
		}).Errorln("[Handler][HandleAuthProviderCallback][Exchange]: ", err.Error())
		return
	}

	idFromProvider, err := h.ServiceAuth.GetIDFromProvider(provider, token)
	if err != nil {
		log.WithFields(log.Fields{
			"provider":          provider,
			"stateCookie":       stateCookie,
			"stateFromCookie":   stateFromCookie,
			"stateFromProvider": stateFromProvider,
			"code":              code,
		}).Errorln("[Handler][HandleAuthProviderCallback][GetIDFromProvider]: ", err.Error())
		return
	}

	_ = idFromProvider
	// TODO: register if idFromProvider is not yet exist, login if idFromProvider
	// already exist

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
