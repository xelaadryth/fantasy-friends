package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xelaadryth/fantasy-friends/database"
)

const (
	loginAction    = "login"
	registerAction = "register"
)

//UserForm input fields required for login and registration
type UserForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Action   string `form:"action" binding:"required"`
}

func login(username string, password string) error {
	err := database.Login(username, password)

	//Obscure the error
	if err != nil {
		return errors.New("Error logging in; most likely invalid credentials.")
	}

	return nil
}

func register(username string, password string) error {
	if username == "" {
		return errors.New("Please enter a username.")
	}
	if password == "" {
		return errors.New("Please enter a password.")
	}

	err := database.Register(username, password)

	//Obscure the error
	if err != nil {
		return errors.New("Error in registration; most likely that username is taken.")
	}

	return nil
}

//processUser logins or registration
func processUser(c *gin.Context) {
	var userForm UserForm
	c.Bind(&userForm)

	err := errors.New("Invalid user action.")

	if userForm.Action == loginAction {
		err = login(userForm.Username, userForm.Password)

	} else if userForm.Action == registerAction {
		err = register(userForm.Username, userForm.Password)
	}

	if err != nil {
		invalidHandler(c, http.StatusBadRequest, err)
	}

	c.Redirect(http.StatusFound, "/")
}
