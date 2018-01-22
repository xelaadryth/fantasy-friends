package fantasy

import (
	"errors"
	"fmt"
)

//Scrim runs a Fantasy match between the summoner names provided
func Scrim(region string, summonerNames ...string) (*MatchScore, error) {
	if len(summonerNames) != PlayersPerMatch {
		err := errors.New(fmt.Sprint("Provided ", len(summonerNames), "summoners instead of ", PlayersPerMatch, "."))
		return &MatchScore{}, err
	}
	summoners, _, err := SummonersBySummonerName(region, summonerNames...)
	if err != nil {
		return &MatchScore{}, err
	}

	return CalculateMatchScore(region, summoners)
}
