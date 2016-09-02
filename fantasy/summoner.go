package fantasy

import (
	"errors"
	"log"

	"github.com/TrevorSStone/goriot"
	"github.com/xelaadryth/fantasy-friends/database"
)

//NormalizeName using goriot to do just one name
func NormalizeName(summonerName string) string {
	summonerNames := make([]string, 1, 1)
	summonerNames[0] = summonerName

	normalizedSummonerNames := goriot.NormalizeSummonerName(summonerNames...)

	return normalizedSummonerNames[0]
}

//attemptUncacheSummonersByName tries to retrieve all specified summoners from cache if possible
func attemptUncacheSummonersByName(region string, summonerNames ...string) (map[string]goriot.Summoner, error) {
	summonersMap := make(map[string]goriot.Summoner)
	var err error

	for _, summonerName := range summonerNames {
		summonersMap[summonerName], err = database.UncacheSummonerByName(region, summonerName)
		if err != nil {
			return summonersMap, err
		}
	}

	//This happens if we have duplicate summoner names
	if len(summonersMap) != len(summonerNames) {
		return summonersMap, errors.New("Duplicate or invalid summoner names.")
	}

	return summonersMap, nil
}

//attemptUncacheSummonersByID tries to retrieve all specified summoner from cache if possible
func attemptUncacheSummonersByID(region string, summonerIDs ...int64) (map[int64]goriot.Summoner, error) {
	summonersMap := make(map[int64]goriot.Summoner)
	var err error

	for _, summonerID := range summonerIDs {
		summonersMap[summonerID], err = database.UncacheSummonerByID(region, summonerID)
		if err != nil {
			return summonersMap, err
		}
	}

	//This happens if we have duplicate summoner names
	if len(summonersMap) != len(summonerIDs) {
		return summonersMap, errors.New("Duplicate or invalid summoner names.")
	}

	return summonersMap, nil
}

//GetSummonersByName by checking the cache and then by asking Riot
//TODO: Pass stuff by reference instead of copying values
func GetSummonersByName(region string, summonerNames ...string) (map[string]goriot.Summoner, error) {
	//Attempt to get Summoner objects for each normalized summoner name
	//TODO: Remove this caching when rate limits are removed
	summonersMap, err := attemptUncacheSummonersByName(region, summonerNames...)
	if err != nil {
		//If we can't, then query the Riot API for the summoner IDs
		summonersMap, err = goriot.SummonerByName(region, summonerNames...)
		if err != nil {
			return summonersMap, err
		}

		//Insert the results of the fresh queries into the summoner cache
		for normalizedName, summoner := range summonersMap {
			cacheErr := database.CacheSummoner(region, normalizedName, summoner)
			if cacheErr != nil {
				log.Println("Failed to cache summoner", normalizedName, summoner, "with error", cacheErr)
			}
		}
	}

	//This happens if we have duplicate summoner names
	if len(summonersMap) != len(summonerNames) {
		return summonersMap, errors.New("Duplicate or invalid summoner names.")
	}

	return summonersMap, nil
}

//GetSummonersByID by checking the cache and then by asking Riot
func GetSummonersByID(region string, summonerIDs ...int64) (map[int64]goriot.Summoner, error) {
	//Attempt to get Summoner objects for each summonerID
	//TODO: Remove this caching when rate limits are removed
	summonersMap, err := attemptUncacheSummonersByID(region, summonerIDs...)
	if err != nil {
		//If we can't, then query the Riot API for the summoner IDs
		summonersMap, err = goriot.SummonerByID(region, summonerIDs...)
		if err != nil {
			return summonersMap, err
		}

		//Insert the results of the fresh queries into the summoner cache
		for _, summoner := range summonersMap {
			normalizedName := NormalizeName(summoner.Name)
			cacheErr := database.CacheSummoner(region, normalizedName, summoner)
			if cacheErr != nil {
				log.Println("Failed to cache summoner", normalizedName, summoner, "with error", cacheErr)
			}
		}
	}

	//This happens if we have duplicate summoner names
	if len(summonersMap) != len(summonerIDs) {
		return summonersMap, errors.New("Duplicate or invalid summoner names.")
	}

	return summonersMap, nil
}
