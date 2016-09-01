package database

// prepareUserStatements readies SQL queries for later use
import (
	"github.com/TrevorSStone/goriot"
	"github.com/jackc/pgx"
)

//Queries that are prepared for easy calling
const (
	QueryGetSummoner    = "getSummoner"
	QueryUpsertSummoner = "upsertSummoner"
)

//UncacheSummoner from database if it exists, return an error otherwise
func UncacheSummoner(region string, normalizedName string) (goriot.Summoner, error) {
	summoner := goriot.Summoner{}
	err := DBConnectionPool.QueryRow(
		QueryGetSummoner, normalizedName, region).Scan(
		&summoner.ID, &summoner.Name, &summoner.SummonerLevel, &summoner.ProfileIconID, &summoner.RevisionDate)

	return summoner, err
}

//CacheSummoner attempts to cache a summoner name if it hasn't been cached before
func CacheSummoner(region string, normalizedName string, summoner goriot.Summoner) error {
	_, err := DBConnectionPool.Exec(QueryUpsertSummoner, summoner.ID, summoner.Name, summoner.SummonerLevel,
		summoner.ProfileIconID, summoner.RevisionDate, normalizedName, region)

	return err
}

func prepareRiotStatements(conn *pgx.Conn) (err error) {
	_, err = conn.Prepare(QueryGetSummoner, `
    SELECT id, summoner_name, summoner_level, profile_icon_id, revision_date
		FROM fantasy_friends.summoner_cache
		WHERE normalized_name=$1 AND region=$2
  `)
	if err != nil {
		return err
	}

	//TODO: Remove this caching when rate limits are removed on production key
	//Edge case when someone name changes (makes us not able to mark other columns as unique)
	//Also has the issue where if someone switches to a an unused name already used by the db, cached version will be used
	_, err = conn.Prepare(QueryUpsertSummoner, `
	  INSERT INTO fantasy_friends.summoner_cache (id, summoner_name, summoner_level, profile_icon_id, revision_date, normalized_name, region)
	  VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (normalized_name)
		DO UPDATE SET (id, summoner_name, summoner_level, profile_icon_id, revision_date, region) = ($1, $2, $3, $4, $5, $7)
	`)

	return err
}
