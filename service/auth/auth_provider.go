package auth

import (
	"io/ioutil"

	"github.com/pikomonde/pikopos/common"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// GenerateStateAndGetAuthURL is used to generate state and get Auth URL
func (s *ServiceAuth) GenerateStateAndGetAuthURL(provider string) (string, string, error) {
	state := common.RandomBase64()

	config, err := s.RepositoryAuth.GetAuthConfig(provider)
	if err != nil {
		log.WithFields(log.Fields{
			"provider": provider,
		}).Errorln("[ServiceAuth][GenerateStateAndGetAuthURL][GetAuthConfig]: ", err.Error())
		return "", "", err
	}

	return state, config.AuthCodeURL(state), nil
}

// Exchange is used to get token from provider using provider's code
func (s *ServiceAuth) Exchange(provider, code string) (*oauth2.Token, error) {
	config, err := s.RepositoryAuth.GetAuthConfig(provider)
	if err != nil {
		log.WithFields(log.Fields{
			"provider": provider,
			"code":     code,
		}).Errorln("[ServiceAuth][Exchange][GetAuthConfig]: ", err.Error())
		return nil, err
	}

	// getting token
	ctx, cancel := common.ContextWithDuration()
	defer cancel()
	token, err := config.Exchange(ctx, code)
	if err != nil {
		log.WithFields(log.Fields{
			"provider": provider,
			"code":     code,
		}).Errorln("[ServiceAuth][Exchange][Exchange]: ", err.Error())
		return nil, err
	}

	return token, nil
}

// GetIDFromProvider is used to get ID from provider.
func (s *ServiceAuth) GetIDFromProvider(provider string, token *oauth2.Token) (string, error) {
	config, err := s.RepositoryAuth.GetAuthConfig(provider)
	if err != nil {
		log.WithFields(log.Fields{
			"provider": provider,
		}).Errorln("[ServiceAuth][GetIDFromProvider][GetAuthConfig]: ", err.Error())
		return "", err
	}

	// getting client
	ctx, cancel := common.ContextWithDuration()
	defer cancel()
	client := config.Client(ctx, token)

	// getting userinfo
	resp, err := client.Get(providers[provider].userinfoURL)
	if err != nil {
		log.WithFields(log.Fields{
			"provider":    provider,
			"userinfoURL": providers[provider].userinfoURL,
		}).Errorln("[ServiceAuth][GetIDFromProvider][Get UserInfo]: ", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	// unmarshall response
	userinfoBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"provider":    provider,
			"userinfoURL": providers[provider].userinfoURL,
		}).Errorln("[ServiceAuth][GetIDFromProvider][ReadAll]: ", err.Error())
		return "", err
	}

	idFromProvider, err := providers[provider].respToID(userinfoBytes)
	if err != nil {
		log.WithFields(log.Fields{
			"provider":    provider,
			"userinfoURL": providers[provider].userinfoURL,
		}).Errorln("[ServiceAuth][GetIDFromProvider][respToID]: ", err.Error())
		return "", err
	}

	return idFromProvider, nil
}
