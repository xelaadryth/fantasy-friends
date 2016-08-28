package database

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

//BcryptCost increases the number of rounds and time taken by the hash algorithm; safe to change this since
//the number of rounds is stored in the bcrypt hash, so will only apply to new passwords
const BcryptCost = 13

//pepper is a fixed series of random characters that is concatenated to the end of passwords before hashing
var pepper string

//Login a user account from the db
func Login(username string, password string) error {
	var id int64
	var hash []byte
	err := DBConnectionPool.QueryRow(
		QueryGetUserAccountByUsername, username).Scan(&id, &hash)
	if err != nil {
		return errors.New("User doesn't exist.")
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password+pepper))
	if err != nil {
		return errors.New("Invalid password.")
	}

	return nil
}

//Register a new account to the db
func Register(username string, password string) error {
	rows, err := DBConnectionPool.Query(QueryGetUserAccountByUsername, username)
	if err != nil {
		return errors.New("Error accessing database.")
	}

	//Next returns false when rows runs out, but we expect to have 0 rows returned so error
	if rows.Next() {
		return errors.New(fmt.Sprint("User ", username, " already exists."))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password+pepper), BcryptCost)
	if err != nil {
		return errors.New("Error generating account.")
	}

	_, err = DBConnectionPool.Exec(QueryPutUserAccount, username, hash)
	if err != nil {
		return errors.New("Error creating account.")
	}

	return nil
}

//PreparePepper gets the pepper value from environment variables
func PreparePepper() error {
	pepper = os.Getenv("PEPPER")
	if pepper == "" {
		return errors.New("$PEPPER must be set")
	}

	return nil
}
