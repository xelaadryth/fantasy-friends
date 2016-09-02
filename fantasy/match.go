package fantasy

import (
	"errors"
	"fmt"

	"github.com/TrevorSStone/goriot"
)

//PlayMatch runs a Fantasy match between the summoners provided
func PlayMatch(region string, inputSummonerNames ...string) (MatchScore, error) {
	if len(inputSummonerNames) != PlayersPerMatch {
		err := errors.New(fmt.Sprint("Provided ", len(inputSummonerNames), "summoners instead of ", PlayersPerMatch, "."))
		return MatchScore{}, err
	}
	summonerNames := goriot.NormalizeSummonerName(inputSummonerNames...)

	summonersMap, _, err := GetSummonersByName(region, summonerNames...)
	if err != nil {
		return MatchScore{}, err
	}

	//Make a data structure with all the player IDs
	summoners := make([]goriot.Summoner, PlayersPerMatch, PlayersPerMatch)
	for i := 0; i < PlayersPerMatch; i++ {
		summoners[i] = summonersMap[summonerNames[i]]
	}
	return CalculateMatchScore(region, summoners)
}
