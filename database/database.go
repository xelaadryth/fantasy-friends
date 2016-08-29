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
			id (bigserial - PK)
			username (varchar(32) - Unique)
			hash (varchar(64))
*/

//DBConnectionPool is required for making queries to the DB
var DBConnectionPool *pgx.ConnPool

//MaxDBConnections the db can support
const MaxDBConnections = 5

//DBTimeout in seconds
const DBTimeout = 30 * time.Second

//Queries that are prepared for easy calling
const (
	QueryGetUserAccountByUsername = "getUserAccount"
	QueryPutUserAccount           = "putUserAccount"
)

// afterConnect creates the prepared statements that this application uses
func afterConnect(conn *pgx.Conn) (err error) {
	_, err = conn.Prepare(QueryGetUserAccountByUsername, `
    SELECT id, hash FROM fantasy_friends.user_account WHERE username=$1
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare(QueryPutUserAccount, `
    INSERT INTO fantasy_friends.user_account(username, hash)
    VALUES ($1, $2)
  `)
	return
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
