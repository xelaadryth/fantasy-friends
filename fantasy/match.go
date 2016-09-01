package fantasy

import (
	"errors"
	"fmt"
	"log"

	"github.com/TrevorSStone/goriot"
	"github.com/xelaadryth/fantasy-friends/database"
)

//attemptUncacheSummoners tries to retrieve all summoner IDs from cache if possible
func attemptUncacheSummoners(region string, normalizedSummonerNames ...string) (map[string]goriot.Summoner, error) {
	summoners := make(map[string]goriot.Summoner)
	var err error

	for _, normalizedSummonerName := range normalizedSummonerNames {
		summoners[normalizedSummonerName], err = database.UncacheSummoner(region, normalizedSummonerName)
		if err != nil {
			return summoners, err
		}
	}

	return summoners, err
}

//PlayMatch runs a Fantasy match between the summoners provided
func PlayMatch(region string, inputSummonerNames ...string) (MatchScore, error) {
	summonerNames := goriot.NormalizeSummonerName(inputSummonerNames...)
	if len(summonerNames) != PlayersPerMatch {
		err := errors.New(fmt.Sprint("Provided ", len(summonerNames), "summoners instead of ", PlayersPerMatch, "."))
		return MatchScore{}, err
	}

	//Attempt to get Summoner objects for each normalized summoner name
	//TODO: Remove this caching when rate limits are removed
	summonersMap, err := attemptUncacheSummoners(region, summonerNames...)
	if err != nil {
		//If we can't, then query the Riot API for the summoner IDs
		summonersMap, err = goriot.SummonerByName(region, summonerNames...)
		if err != nil {
			return MatchScore{}, err
		}

		//Insert the results of the fresh queries into the summoner cache
		for normalizedName, summoner := range summonersMap {
			cacheErr := database.CacheSummoner(region, normalizedName, summoner)
			if cacheErr != nil {
				log.Println("Failed to cache summoner", normalizedName, summoner, "with error", cacheErr)
			}
		}
	}
	//Make sure the names are unique
	if len(summonersMap) != PlayersPerMatch {
		err = errors.New(
			fmt.Sprint("Only ", len(summonersMap), " valid distinct summoner names instead of ",
				PlayersPerMatch, "."))
		return MatchScore{}, err
	}

	//Make a data structure with all the player IDs
	summoners := make([]goriot.Summoner, PlayersPerMatch, PlayersPerMatch)
	for i := 0; i < PlayersPerMatch; i++ {
		summoners[i] = summonersMap[summonerNames[i]]
	}
	return CalculateMatchScore(region, summoners)
}
