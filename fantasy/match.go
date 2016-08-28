package fantasy

import (
	"errors"
	"fmt"

	"github.com/TrevorSStone/goriot"
)

//PlayMatch runs a Fantasy match between the summoners provided
func PlayMatch(region string, inputSummonerNames ...string) (MatchScore, error) {
	summonerNames := goriot.NormalizeSummonerName(inputSummonerNames...)
	if len(summonerNames) != PlayersPerMatch {
		return MatchScore{}, errors.New(
			fmt.Sprint("Provided ", len(summonerNames), "summoners instead of ", PlayersPerMatch, "."))
	}

	//Attempt to get the summoner IDs
	summonersMap, err := goriot.SummonerByName(region, summonerNames...)
	if err != nil {
		return MatchScore{}, err
	}
	//Make sure the names are unique
	if len(summonersMap) != PlayersPerMatch {
		return MatchScore{}, errors.New(
			fmt.Sprint("Only ", len(summonersMap), " valid distinct summoner names instead of ",
				PlayersPerMatch, "."))
	}

	//Make a data structure with all the player IDs
	summoners := make([]goriot.Summoner, PlayersPerMatch, PlayersPerMatch)
	for i := 0; i < PlayersPerMatch; i++ {
		summoners[i] = summonersMap[summonerNames[i]]
	}
	matchScore, err := CalculateScores(region, summoners)
	if err != nil {
		return MatchScore{}, err
	}

	return matchScore, nil
}
