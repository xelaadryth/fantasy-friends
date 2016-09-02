package database

// prepareUserStatements readies SQL queries for later use
import (
	"github.com/TrevorSStone/goriot"
	"github.com/jackc/pgx"
)

//Queries that are prepared for easy calling
const (
	QueryGetSummonerByName       = "getSummonerByName"
	QueryGetSummonerBySummonerID = "getSummonerBySummonerID"
	QueryGetSummonerByCacheID    = "getSummonerByCacheID"
	QueryUpsertSummoner          = "upsertSummoner"
)

//UncacheSummonerByName from database if it exists, return an error otherwise
//TODO: Do this in bulk, and perform the Riot API queries/caching here
func UncacheSummonerByName(region string, normalizedName string) (goriot.Summoner, int64, error) {
	var cacheID int64
	summoner := goriot.Summoner{}
	err := DBConnectionPool.QueryRow(
		QueryGetSummonerByName, normalizedName, region).Scan(
		&cacheID, &summoner.ID, &summoner.Name, &summoner.SummonerLevel, &summoner.ProfileIconID, &summoner.RevisionDate)

	return summoner, cacheID, err
}

//UncacheSummonerBySummonerID from database if it exists, return an error otherwise
//TODO: Do this in bulk, and perform the Riot API queries/caching here
func UncacheSummonerBySummonerID(region string, summonerID int64) (goriot.Summoner, int64, error) {
	var cacheID int64
	summoner := goriot.Summoner{}
	err := DBConnectionPool.QueryRow(
		QueryGetSummonerBySummonerID, summonerID, region).Scan(
		&cacheID, &summoner.ID, &summoner.Name, &summoner.SummonerLevel, &summoner.ProfileIconID, &summoner.RevisionDate)

	return summoner, cacheID, err
}

//UncacheSummonerByCacheID from database if it exists, return an error otherwise
func UncacheSummonerByCacheID(summonerCacheID int64) (goriot.Summoner, error) {
	summoner := goriot.Summoner{}
	err := DBConnectionPool.QueryRow(
		QueryGetSummonerByCacheID, summonerCacheID).Scan(
		&summoner.ID, &summoner.Name, &summoner.SummonerLevel, &summoner.ProfileIconID, &summoner.RevisionDate)

	//TODO: If it's too old then delete this entry so it will be refreshed next time

	return summoner, err
}

//CacheSummoner attempts to cache a summoner if it hasn't been cached before
func CacheSummoner(region string, normalizedName string, summoner goriot.Summoner) (int64, error) {
	var summonerCacheID int64
	err := DBConnectionPool.QueryRow(QueryUpsertSummoner, summoner.ID, summoner.Name, summoner.SummonerLevel,
		summoner.ProfileIconID, summoner.RevisionDate, normalizedName, region).Scan(&summonerCacheID)

	return summonerCacheID, err
}

func prepareRiotStatements(conn *pgx.Conn) error {
	_, err := conn.Prepare(QueryGetSummonerByName, `
		SELECT id, summoner_id, summoner_name, summoner_level, profile_icon_id, revision_date
		FROM fantasy_friends.summoner_cache
		WHERE normalized_name=$1 AND region=$2
  `)
	if err != nil {
		return err
	}

	_, err = conn.Prepare(QueryGetSummonerBySummonerID, `
			SELECT id, summoner_id, summoner_name, summoner_level, profile_icon_id, revision_date
			FROM fantasy_friends.summoner_cache
			WHERE summoner_id=$1 AND region=$2
	  `)
	if err != nil {
		return err
	}

	_, err = conn.Prepare(QueryGetSummonerByCacheID, `
			SELECT summoner_id, summoner_name, summoner_level, profile_icon_id, revision_date
			FROM fantasy_friends.summoner_cache
			WHERE id=$1
	  `)
	if err != nil {
		return err
	}

	_, err = conn.Prepare(QueryUpsertSummoner, `
		INSERT INTO fantasy_friends.summoner_cache (summoner_id, summoner_name, summoner_level, profile_icon_id, revision_date, normalized_name, region)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (summoner_id, region)
		DO UPDATE SET (summoner_name, summoner_level, profile_icon_id, revision_date, normalized_name) = ($2, $3, $4, $5, $6)
		RETURNING id
	`)

	return err
}
