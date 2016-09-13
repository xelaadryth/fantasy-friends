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

//GetSummonersByName from cache if available, otherwise grab from Riot API and cache them.
//Outputs a mapping of normalized summoner name to summoner objects
func GetSummonersByName(region string, inputSummonerNames ...string) (map[string]goriot.Summoner, map[string]int64, error) {
	summonerNames := goriot.NormalizeSummonerName(inputSummonerNames...)
	nameToSummoner := make(map[string]goriot.Summoner)
	nameToCacheID := make(map[string]int64)
	var err error
	//Try to uncache them first
	for _, summonerName := range summonerNames {
		nameToSummoner[summonerName], nameToCacheID[summonerName], err = database.UncacheSummonerByName(region, summonerName)
		if err != nil {
			break
		}
	}

	//Cache contained all of the summoner names, we're done
	if len(nameToSummoner) == len(inputSummonerNames) && err == nil {
		return nameToSummoner, nameToCacheID, nil
	}

	//TODO: Move this into database in summoner_cache
	//Either the cache didn't contain all of the ids or there are repeat IDs. Query Riot API next
	nameToSummoner, err = goriot.SummonerByName(region, summonerNames...)
	if err != nil || len(nameToSummoner) != len(inputSummonerNames) {
		return nameToSummoner, nameToCacheID, errors.New("Duplicate or invalid summoner names.")
	}

	var cacheErr error
	//Insert the results of the fresh query into the summoner cache
	for normalizedName, summoner := range nameToSummoner {
		nameToCacheID[normalizedName], cacheErr = database.CacheSummoner(region, normalizedName, summoner)
		if cacheErr != nil {
			log.Println("Failed to cache summoner", normalizedName, summoner, "with error", cacheErr)
		}
	}

	return nameToSummoner, nameToCacheID, nil
}

//GetSummonersBySummonerID from cache if available, otherwise grab from Riot API and cache them
func GetSummonersBySummonerID(region string, summonerIDs ...int64) (map[int64]goriot.Summoner, map[int64]int64, error) {
	idToSummoner := make(map[int64]goriot.Summoner)
	idToCacheID := make(map[int64]int64)
	var err error
	//Try to uncache them first
	for _, summonerID := range summonerIDs {
		idToSummoner[summonerID], idToCacheID[summonerID], err = database.UncacheSummonerBySummonerID(region, summonerID)
		if err != nil {
			break
		}
	}

	//Cache contained all of the summoner IDs, we're done
	if len(idToSummoner) == len(summonerIDs) && err == nil {
		return idToSummoner, idToCacheID, nil
	}

	//TODO: Move this into database in summoner_cache
	//Either the cache didn't contain all of the ids or there are repeat IDs. Query Riot API next
	idToSummoner, err = goriot.SummonerByID(region, summonerIDs...)
	if err != nil || len(idToSummoner) != len(summonerIDs) {
		return idToSummoner, idToCacheID, errors.New("Duplicate or invalid summoner IDs.")
	}

	var cacheErr error
	//Insert the results of the fresh query into the summoner cache
	for _, summoner := range idToSummoner {
		normalizedName := NormalizeName(summoner.Name)
		idToCacheID[summoner.ID], cacheErr = database.CacheSummoner(region, normalizedName, summoner)
		if cacheErr != nil {
			log.Println("Failed to cache summoner", normalizedName, summoner, "with error", cacheErr)
		}
	}

	return idToSummoner, idToCacheID, nil
}

//GetSummonersByCacheID if they exist, error otherwise. Should never error due to foreign key constraints unless
//malicious users are messing around
func GetSummonersByCacheID(summonerCacheIDs ...int64) (map[int64]goriot.Summoner, error) {
	cacheIDToSummoner := make(map[int64]goriot.Summoner)
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
