package riot

import (
	"fmt"
	"log"
	"strings"
)

const summonerEndpoint = "/api/lol/%s/v1.4/summoner/by-name/%s"

//Summoner matches the v1.4 endpoint from Riot API
type Summoner struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

//standardizeName removes spaces and makes a name lowercase
func standardizeName(summonerName string) string {
	return strings.ToLower(strings.Replace(summonerName, " ", "", -1))
}

//GetSummoner gets the summoner object for a given summoner name
func GetSummoner(region string, summonerName string) (*Summoner, error) {
	//JSON form of response is a list of summoners that match the query
	results := make(map[string]Summoner)
	err := pathGet(fmt.Sprintf(summonerEndpoint, region, standardizeName(summonerName)), &results)

	var currentSummoner Summoner
	for _, value := range results {
		currentSummoner = value
		break
	}

	return &currentSummoner, err
}

//GetSummonerID gets the summoner ID for a given summoner name
func GetSummonerID(region string, summonerName string) (int, error) {
	s, err := GetSummoner(region, summonerName)
	if err != nil {
		//TODO: Error handling
		log.Println(err)

		return 0, err
	}
	return s.ID, err
}
