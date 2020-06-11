package clients

import (
	"fmt"

	// initialize mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pikomonde/pikopos/config"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
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
	for provider := range config.C.OAuth {
		// set endpoint
		endpoint := oauth2.Endpoint{}
		switch provider {
		case "google":
			endpoint = endpoints.Google
		case "facebook":
			endpoint = endpoints.Facebook
		}

		// set oauthconfig
		oauthConfig[provider] = &oauth2.Config{
			RedirectURL:  fmt.Sprintf("https://%s/auth/%s/callback", config.C.BaseURL, provider),
			ClientID:     config.C.OAuth[provider].ClientID,
			ClientSecret: config.C.OAuth[provider].ClientSecret,
			Scopes:       config.C.OAuth[provider].Scopes,
			Endpoint:     endpoint,
		}
	}

	return oauthConfig
}
