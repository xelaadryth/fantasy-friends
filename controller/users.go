package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xelaadryth/fantasy-friends/database"
	"golang.org/x/crypto/bcrypt"
)

//BcryptCost increases the number of rounds and time taken by the hash algorithm
const BcryptCost = 12

//UserForm input fields required for login and registration
type UserForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func login(c *gin.Context) {
	var userForm UserForm
	c.Bind(&userForm)

	var id int64
	var hash []byte
	err := database.DBConnectionPool.QueryRow(
		database.QueryGetUserAccountByUsername, userForm.Username).Scan(&id, &hash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Invalid login credentials.",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(userForm.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Invalid login credentials.",
		})
		return
	}

	log.Println("Login for", userForm.Username, "successful!")

	c.Redirect(http.StatusFound, "/")
}

func register(c *gin.Context) {
	var userForm UserForm
	c.Bind(&userForm)

	rows, err := database.DBConnectionPool.Query(database.QueryGetUserAccountByUsername, userForm.Username)
	if err != nil {
		log.Println("Error on query to user accounts table.")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Error accessing database.",
		})
		return
	}

	//Next returns false when rows runs out, but we expect to have 0 rows returned so error
	if rows.Next() {
		rows.Close()
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  fmt.Sprint("User ", userForm.Username, " already exists."),
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(userForm.Password), BcryptCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Error generating account.",
		})
		return
	}

	_, err = database.DBConnectionPool.Exec(database.QueryPutUserAccount, userForm.Username, hash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Error in registration, please try again.",
		})
		return
	}

	log.Println("Registration for", userForm.Username, "successful!")

	c.Redirect(http.StatusFound, "/")
}
