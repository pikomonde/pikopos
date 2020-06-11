package auth

import "encoding/json"

var providers = map[string]provider{
	"google": provider{
		userinfoURL: "https://www.googleapis.com/oauth2/v3/userinfo",
		respToID:    googleRespToUserInfo,
	},
	"facebook": provider{
		// userinfoURL: "https://graph.facebook.com/me?fields=email,first_name,last_name,link,about,id,name,picture,location",
		userinfoURL: "https://graph.facebook.com/me?fields=id",
		respToID:    facebookRespToUserInfo,
	},
}

type provider struct {
	userinfoURL string
	respToID    func([]byte) (string, error)
}

func googleRespToUserInfo(respByte []byte) (string, error) {
	resp := struct {
		ID string `json:"sub"`
	}{}
	err := json.Unmarshal(respByte, &resp)
	return resp.ID, err
}

func facebookRespToUserInfo(respByte []byte) (string, error) {
	resp := struct {
		ID string `json:"id"`
	}{}
	err := json.Unmarshal(respByte, &resp)
	return resp.ID, err
}
