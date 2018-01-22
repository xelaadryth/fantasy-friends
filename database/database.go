package database

import (
	"crypto/tls"
	"os"
	"time"

	"github.com/jackc/pgx"
)

/*
Currently the PostgreSQL database structure looks like so (in this order):

database_name (Database)
	fantasy_friends (Schema)
		user_account (Table)
			id (bigserial - PK, NN)
			username (text - Unique, NN)
			username_lower (text - Unique, NN)
			hashed_password (text - NN)
		user_session (Table)
			id (text - PK, NN)
			user_id (bigint - NN, FK(user_account.id))
			creation_time (bigint - NN)
		fantasy_players (Table)
			id (bigserial - PK, NN)
			account_id (bigint - NN, Unique with region)
			summoner_id (bignit - NN, Unique with region)
			summoner_name (text - NN)
			summoner_level (int - NN)
			profile_icon_id (int - NN)
			revision_date (bigint - NN)
			normalized_name (text - NN)
			region (text - NN)
		fantasy_team (Table)
			id (bigserial - PK, NN)
			owner (bigint - NN, FK(user_account.id))
			team_name (text - NN)
			region (text - NN)
			position (int - NN, Unique with owner)
			top (bigint - FK(fantasy_player.id))
			jungle (bigint - FK(fantasy_player.id))
			mid (bigint - FK(fantasy_player.id))
			bottom (bigint - FK(fantasy_player.id))
			support (bigint - FK(fantasy_player.id))
*/

//DBConnectionPool is required for making queries to the DB
var DBConnectionPool *pgx.ConnPool

//MaxDBConnections the db can support
const MaxDBConnections = 5

//DBTimeout in seconds
const DBTimeout = 30 * time.Second

// afterConnect creates the prepared statements that this application uses
func afterConnect(conn *pgx.Conn) error {
	err := prepareUserStatements(conn)
	if err != nil {
		return err
	}

	err = prepareSessionStatements(conn)
	if err != nil {
		return err
	}

	err = preparePlayerStatements(conn)
	if err != nil {
		return err
	}

	err = prepareTeamStatements(conn)
	if err != nil {
		return err
	}

	return nil
}

//Connect to the database
func Connect() error {

	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DB"),
			TLSConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		MaxConnections: MaxDBConnections,
		AfterConnect:   afterConnect,
		AcquireTimeout: DBTimeout,
	}
	pool, err := pgx.NewConnPool(connPoolConfig)
	if err != nil {
		return err
	}
	DBConnectionPool = pool

	return nil
}
