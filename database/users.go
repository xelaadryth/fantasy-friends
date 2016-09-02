package database

import (
	"errors"
	"os"

	"github.com/jackc/pgx"
	"golang.org/x/crypto/bcrypt"
)

//BcryptCost increases the number of rounds and time taken by the hash algorithm; safe to change this since
//the number of rounds is stored in the bcrypt hash, so will only apply to new passwords
const BcryptCost = 13

//pepper is a fixed series of random characters that is concatenated to the end of passwords before hashing
var pepper string

//Queries that are prepared for easy calling
const (
	QueryGetUserAccountByUsername = "getUserAccount"
	QueryPutUserAccount           = "putUserAccount"
)

//Login a user account from the db
func Login(username string, password string) (string, error) {
	var userID int64
	var hash []byte
	err := DBConnectionPool.QueryRow(
		QueryGetUserAccountByUsername, username).Scan(&userID, &hash)
	if err != nil {
		return "", errors.New("User doesn't exist.")
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password+pepper))
	if err != nil {
		return "", errors.New("Invalid password.")
	}

	sessionID, err := AddSession(userID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

//Register a new account to the db
func Register(username string, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password+pepper), BcryptCost)
	if err != nil {
		return "", errors.New("Error generating account.")
	}

	var userID int64
	err = DBConnectionPool.QueryRow(QueryPutUserAccount, username, hash).Scan(&userID)
	if err != nil {
		return "", errors.New("Error creating account.")
	}
	sessionID, err := AddSession(userID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

//PreparePepper gets the pepper value from environment variables
func PreparePepper() error {
	pepper = os.Getenv("PEPPER")
	if pepper == "" {
		return errors.New("$PEPPER must be set (to anything)")
	}

	return nil
}

// prepareUserStatements readies SQL queries for later use
func prepareUserStatements(conn *pgx.Conn) error {
	_, err := conn.Prepare(QueryGetUserAccountByUsername, `
		SELECT id, hash FROM fantasy_friends.user_account WHERE username=$1
  `)
	if err != nil {
		return err
	}

	_, err = conn.Prepare(QueryPutUserAccount, `
		INSERT INTO fantasy_friends.user_account (username, hash)
		VALUES ($1, $2)
		RETURNING id
  `)
	if err != nil {
		return err
	}

	return nil
}
