package auth

import (
	"github.com/pikomonde/pikopos/common"
	log "github.com/sirupsen/logrus"
)

// GetAuthStateAndURL is used to get the oauth state and code URL
func (s *ServiceAuth) GetAuthStateAndURL(provider string) (string, string, error) {
	state := common.RandomBase64()

	config, err := s.RepositoryAuth.GetAuthConfig(provider)
	if err != nil {
		log.WithFields(log.Fields{
			"provider": provider,
		}).Errorln("[ServiceAuth][GetAuthURL][GetAuthConfig]: ", err.Error())
		return "", "", err
	}

	return state, config.AuthCodeURL(state), nil
}
