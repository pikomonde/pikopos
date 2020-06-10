package clients

import (
	"fmt"

	// initialize mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pikomonde/pikopos/common"
	"github.com/pikomonde/pikopos/config"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Clients contains all clients that will be used in the repository
type Clients struct {
	PikoposMySQLCli *sqlx.DB
	OAuthConfig     map[string]*oauth2.Config
}

// New getting all clients that will be used in the repository
func New() *Clients {
	// Initialize pikopos MySQL
	log.Infoln("[Clients]: Initialize pikopos MySQL...")
	pikoposMySQLCli, err := newPikoposMySQL()
	if err != nil {
		log.Panicln("[Clients]: Failed initialize pikopos MySQL", err.Error())
	}

	// Initialize OAuth config
	log.Infoln("[Clients]: Initialize oauth...")
	oauthConf := newOAuth()

	return &Clients{
		PikoposMySQLCli: pikoposMySQLCli,
		OAuthConfig:     oauthConf,
	}
}

func newPikoposMySQL() (*sqlx.DB, error) {
	return sqlx.Connect("mysql", config.C.MySQLPikopos)
}

func newOAuth() map[string]*oauth2.Config {
	oauthConfig := make(map[string]*oauth2.Config)
	oauthConfig["google"] = &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/auth/google/callback", config.C.BaseURL),
		ClientID:     config.C.OAuthGoogle.ClientID,
		ClientSecret: config.C.OAuthGoogle.ClientSecret,
		Scopes: []string{
			common.OAuthGoogleEmailScope,
		},
		Endpoint: google.Endpoint,
	}

	return oauthConfig

	// fmt.Println(googleOauthConfig.AuthCodeURL("AMSMASMAMSAMSMAMSAMS"))
	// token, _ := googleOauthConfig.Exchange(oauth2.NoContext, "")
	// client := googleOauthConfig.Client(oauth2.NoContext, token)
	// client.Get("https://www.googleapis.com/auth/userinfo.email")

	// return &http.Client{}, nil
}
