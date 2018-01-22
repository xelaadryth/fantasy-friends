package fantasy

import (
	"errors"
	"log"

	"github.com/xelaadryth/fantasy-friends/database"
	"github.com/xelaadryth/fantasy-friends/rgapi"
)

//AccountIdsBySummonerName from cache if available, otherwise grab from Riot API and cache them.
//Outputs a mapping of normalized summoner name to summoner objects
func SummonersBySummonerName(region string, summonerNames ...string) ([]rgapi.Summoner, map[string]int64, error) {
	summoners := make([]rgapi.Summoner, len(summonerNames))
	cacheIds := make(map[string]int64)

	for i, summonerName := range summonerNames {
		//TODO: Save players with additional features, instead of in a cache
		normalizedName := rgapi.NormalizeGameName(summonerName)

		//Check for a cached value
		summoner, cacheId, err := database.UncacheSummonerByName(region, normalizedName)
		if err == nil {
			summoners[i] = summoner
			cacheIds[normalizedName] = cacheId
			continue
		}

		//Not cached, call rgapi
		summoner, err = rgapi.SummonerByName(region, normalizedName)
		if err != nil {
			return summoners, cacheIds, err
		}

		summoners[i] = summoner

		//Cache the value
		cacheIds[normalizedName], err = database.CacheSummoner(region, normalizedName, summoner)
		if err != nil {
			log.Println("Failed to cache summoner", normalizedName, summoner, "with error", err)
		}
	}

	return summoners, cacheIds, nil
}

//GetSummonersByCacheID if they exist, error otherwise. Should never error due to foreign key constraints unless
//malicious users are messing around
func GetSummonersByCacheID(summonerCacheIDs ...int64) (map[int64]rgapi.Summoner, error) {
	cacheIDToSummoner := make(map[int64]rgapi.Summoner)
	var err error
	//Try to uncache them first
	for _, summonerCacheID := range summonerCacheIDs {
		cacheIDToSummoner[summonerCacheID], err = database.UncacheSummonerByCacheID(summonerCacheID)
		if err != nil {
			return cacheIDToSummoner, errors.New("Failed to uncache summoner.")
		}
	}

	//Cache contained all of the summoner IDs, we're done
	if len(cacheIDToSummoner) != len(summonerCacheIDs) {
		return cacheIDToSummoner, errors.New("Duplicate or invalid summoner cache IDs.")
	}

	return cacheIDToSummoner, nil
}
