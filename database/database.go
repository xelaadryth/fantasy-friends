package database

import (
	"os"
	"time"

	"github.com/jackc/pgx"
)

/*
Currently the PostgreSQL database structure looks like so:

database_name (Database)
	fantasy_friends (Schema)
		user_account (Table)
			id (bigserial - PK, NN)
			username (varchar(32) - Unique, NN)
			hash (varchar(64) - NN)
		user_session (Table)
			id (varchar(64) - PK, NN)
			user_id (bigint - NN, FK(user_account.id))
			creation_time (bigint - NN)
		summoner_cache (Table)
			id (bigint - NN)
			summoner_name (varchar(32) - NN)
			summoner_level (int - NN)
			profile_icon_id (int - NN)
			revision_date (bigint - NN)
			normalized_name (varchar(32) - PK, NN)
			region (varchar(8) - NN)
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

	err = prepareRiotStatements(conn)
	if err != nil {
		return err
	}

	err = prepareSessionStatements(conn)
	if err != nil {
		return err
	}

	return nil
}

//Connect to the database
func Connect() error {
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_DATABASE"),
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
