package riot

import (
	"net/http"
	"time"

	"github.com/xelaadryth/fantasy-friends/lib/webhelpers"
)

var urlBase = "https://na.api.pvp.net"
var urlSuffix = "?api_key="
var apiKey = ""

//Client for making get requests
var riotClient = &http.Client{
	Timeout: time.Second * 10,
}

func pathGet(path string, target interface{}) error {
	return webhelpers.GetJSON(urlBase+path+urlSuffix+apiKey, target)
}

//SetAPIKey changes the API key used to connect to the Riot API
func SetAPIKey(key string) {
	apiKey = key
}
