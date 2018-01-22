package database

// prepareUserStatements readies SQL queries for later use
import (
	"strings"

	"github.com/jackc/pgx"

	"github.com/xelaadryth/fantasy-friends/rgapi"
)

//Queries that are prepared for easy calling
const (
	QueryGetSummonerByName      = "getSummonerByName"
	QueryGetSummonerByAccountId = "getSummonerByAccountId"
	QueryGetSummonerByCacheID   = "getSummonerByCacheID"
	QueryUpsertSummoner         = "upsertSummoner"
)

//UncacheSummonerByName from database if it exists, return an error otherwise
//TODO: Do this in bulk, and perform the Riot API queries/caching here
func UncacheSummonerByName(region string, normalizedName string) (rgapi.Summoner, int64, error) {
	var cacheID int64
	summoner := rgapi.Summoner{}
	err := DBConnectionPool.QueryRow(
		QueryGetSummonerByName, normalizedName, strings.ToLower(region)).Scan(
		&cacheID, &summoner.AccountId, &summoner.SummonerId, &summoner.Name, &summoner.SummonerLevel, &summoner.ProfileIconId, &summoner.RevisionDate)

	return summoner, cacheID, err
}

//UncacheSummonerByAccountId from database if it exists, return an error otherwise
//TODO: Do this in bulk, and perform the Riot API queries/caching here
func UncacheSummonerByAccountId(region string, accountId int64) (rgapi.Summoner, int64, error) {
	var cacheID int64
	summoner := rgapi.Summoner{}
	err := DBConnectionPool.QueryRow(
		QueryGetSummonerByAccountId, accountId, strings.ToLower(region)).Scan(
		&cacheID, &summoner.AccountId, &summoner.SummonerId, &summoner.Name, &summoner.SummonerLevel, &summoner.ProfileIconId, &summoner.RevisionDate)

	return summoner, cacheID, err
}

//UncacheSummonerByCacheID from database if it exists, return an error otherwise
func UncacheSummonerByCacheID(summonerCacheID int64) (rgapi.Summoner, error) {
	summoner := rgapi.Summoner{}
	err := DBConnectionPool.QueryRow(
		QueryGetSummonerByCacheID, summonerCacheID).Scan(
		&summoner.AccountId, &summoner.SummonerId, &summoner.Name, &summoner.SummonerLevel, &summoner.ProfileIconId, &summoner.RevisionDate)

	//TODO: If it's too old then delete this entry so it will be refreshed next time

	return summoner, err
}

//CacheSummoner attempts to cache a summoner if it hasn't been cached before
func CacheSummoner(region string, normalizedName string, summoner rgapi.Summoner) (int64, error) {
	var summonerCacheID int64
	err := DBConnectionPool.QueryRow(QueryUpsertSummoner, summoner.AccountId, summoner.SummonerId, summoner.Name, summoner.SummonerLevel,
		summoner.ProfileIconId, summoner.RevisionDate, normalizedName, strings.ToLower(region)).Scan(&summonerCacheID)

	return summonerCacheID, err
}

func preparePlayerStatements(conn *pgx.Conn) error {
	_, err := conn.Prepare(QueryGetSummonerByName, `
		SELECT id, account_id, summoner_id, summoner_name, summoner_level, profile_icon_id, revision_date
		FROM fantasy_friends.fantasy_player
		WHERE normalized_name=$1 AND region=$2
  `)
	if err != nil {
		return err
	}

	_, err = conn.Prepare(QueryGetSummonerByAccountId, `
			SELECT id, account_id, summoner_id, summoner_name, summoner_level, profile_icon_id, revision_date
			FROM fantasy_friends.fantasy_player
			WHERE account_id=$1 AND region=$2
	  `)
	if err != nil {
		return err
	}

	_, err = conn.Prepare(QueryGetSummonerByCacheID, `
			SELECT account_id, summoner_id, summoner_name, summoner_level, profile_icon_id, revision_date
			FROM fantasy_friends.fantasy_player
			WHERE id=$1
	  `)
	if err != nil {
		return err
	}

	_, err = conn.Prepare(QueryUpsertSummoner, `
		INSERT INTO fantasy_friends.fantasy_player (account_id, summoner_id, summoner_name, summoner_level, profile_icon_id, revision_date, normalized_name, region)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (account_id, region)
		DO UPDATE SET (summoner_id, summoner_name, summoner_level, profile_icon_id, revision_date, normalized_name) = ($2, $3, $4, $5, $6, $7)
		RETURNING id
	`)

	return err
}
