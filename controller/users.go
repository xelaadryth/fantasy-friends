package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xelaadryth/fantasy-friends/database"
)

//UserForm input fields required for login and registration
type UserForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func login(c *gin.Context) {
	var userForm UserForm
	c.Bind(&userForm)

	err := database.Login(userForm.Username, userForm.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Error logging in; most likely invalid credentials.",
		})
		return
	}

	log.Println("Login for", userForm.Username, "successful!")

	c.Redirect(http.StatusFound, "/")
}

func register(c *gin.Context) {
	var userForm UserForm
	c.Bind(&userForm)

	if userForm.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Please enter a username.",
		})
		return
	}
	if userForm.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Please enter a password.",
		})
		return
	}

	err := database.Register(userForm.Username, userForm.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Error in registration; most likely that username is taken.",
		})
		return
	}

	log.Println("Registration for", userForm.Username, "successful!")

	c.Redirect(http.StatusFound, "/")
}
