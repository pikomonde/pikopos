package localconfig

import (
	"fmt"

	"github.com/pikomonde/pikopos/clients"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// RepositoryAuth contains clients and Auth repositories
type RepositoryAuth struct {
	Clients *clients.Clients
}

// GetAuthConfig is used to get repository's authentication config
func (r RepositoryAuth) GetAuthConfig(provider string) (*oauth2.Config, error) {
	oauthConfig, ok := r.Clients.OAuthConfig[provider]
	if !ok {
		err := fmt.Errorf("No OAuthConfig exist for the following provider: %s", provider)
		log.WithFields(log.Fields{
			"provider": provider,
		}).Errorln("[RepositoryAuth][GetAuthConfig]: ", err.Error())
		return nil, err
	}
	return oauthConfig, nil
}
