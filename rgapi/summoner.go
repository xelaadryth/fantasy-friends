package rgapi

import (
	"fmt"
	"strings"
)

//Summoner is a player of League of Legends
type Summoner struct {
	ProfileIconId int    `json:"profileIconId"`
	Name          string `json:"name"`
	SummonerLevel int    `json:"summonerLevel"`
	AccountId     int64  `json:"accountId"`
	SummonerId    int64  `json:"id"`
	RevisionDate  int64  `json:"revisionDate"`
}

//NormalizeGameName takes an arbitrary number of strings and returns a string array containing the strings
//standardized to league of legends internal standard (lowecase and strings removed)
func NormalizeGameName(summonerName string) string {
	summonerName = strings.ToLower(summonerName)
	summonerName = strings.Replace(summonerName, " ", "", -1)
	return summonerName
}

//NormalizeGameName takes an arbitrary number of strings and returns a string array containing the strings
//standardized to league of legends internal standard (lowecase and strings removed)
func NormalizeGameNames(summonerNames ...string) []string {
	normalizedNames := make([]string, len(summonerNames))
	for i, summonerName := range summonerNames {
		normalizedNames[i] = NormalizeGameName(summonerName)
	}
	return normalizedNames
}

//SummonerByName retrieves the summoner information of the provided summoner names from Riot Games API.
//It returns a Map of Summoner with the key being the summoner name and any errors that occured from the server
//The global API key must be set before use
//WARNING: The map's key is not necessarily the same string used in the request. It is
//recommended to use NormalizeGameName before calling this function
func SummonerByName(region string, name string) (summoner Summoner, err error) {
	err = apiGet(region, fmt.Sprintf("/summoner/v3/summoners/by-name/%s", strings.ToLower(name)), &summoner)
	return
}
