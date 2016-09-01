package database

import (
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx"
	"github.com/xelaadryth/fantasy-friends/utils"
)

//Queries that are prepared for easy calling
const (
	QueryGetSessionUserID = "getSession"
	QueryPutSession       = "putSession"
)

//InsertionAttempts to generate a session ID. Honestly we should never even have a conflict
const InsertionAttempts = 5
const SessionIDLength = 64

func GetUserIDFromSession(sessionID string) (int64, error) {
	var userID int64
	err := DBConnectionPool.QueryRow(
		QueryGetSessionUserID, sessionID).Scan(&userID)
	if err != nil {
		return 0, errors.New("Invalid session ID.")
	}

	return userID, nil
}

func AddSession(userID int64) (string, error) {
	//Keep trying to insert until we have a unique sessionID
	for i := 0; i < InsertionAttempts; i++ {
		sessionID, err := utils.GenerateString(SessionIDLength)

		if err != nil {
			return sessionID, err
		}

		//There's theoretically a chance that the sessionID was already taken, so try again
		_, err = DBConnectionPool.Exec(QueryPutSession, sessionID, userID, time.Now().Unix())
		if err != nil {
			log.Println(err)
			continue
		}

		return sessionID, nil
	}

	return "", errors.New("Unable to generate a unique session ID.")
}

//prepareSessionStatements readies SQL queries for later use
func prepareSessionStatements(conn *pgx.Conn) error {
	_, err := conn.Prepare(QueryGetSessionUserID, `
    SELECT user_id FROM fantasy_friends.user_session WHERE id=$1
  `)
	if err != nil {
		return err
	}

	_, err = conn.Prepare(QueryPutSession, `
    INSERT INTO fantasy_friends.user_session (id, user_id, creation_time)
    VALUES ($1, $2, $3)
  `)
	if err != nil {
		return err
	}

	//TODO: Add a query to delete all entries that are older than 30 days and run on a cadence

	return nil
}
