package arm

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/golang/glog"
)

// StolenToken is the token generated by azure-cli, i.e. ~/.azure/accessTokens.json
type StolenToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ClientID     string `json:"_clientId"`
}

// StoleToken returns the oauth2 token generated by azure-cli
func StoleToken() *StolenToken {
	accessTokenFile := path.Join(os.Getenv("HOME"), ".azure", "accessTokens.json")
	bs, err := ioutil.ReadFile(accessTokenFile)
	if err != nil {
		glog.Exit(err)
	}
	var toks []StolenToken
	if err := json.Unmarshal(bs, &toks); err != nil {
		glog.Exit(err)
	}
	for _, tk := range toks {
		return &tk
	}
	return nil
}
