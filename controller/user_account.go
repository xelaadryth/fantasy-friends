package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
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

func login(username string, password string) (string, error) {
	sessionID, err := database.Login(username, password)

	//Obscure the error
	if err != nil {
		return "", errors.New("Error logging in; most likely invalid credentials.")
	}

	return sessionID, nil
}

func register(username string, password string) (string, error) {
	if username == "" {
		return "", errors.New("Please enter a username.")
	}
	if password == "" {
		return "", errors.New("Please enter a password.")
	}

	sessionID, err := database.Register(username, password)

	//Obscure the error
	if err != nil {
		return "", errors.New("Error in registration; most likely that username is taken.")
	}

	return sessionID, nil
}

//processUser logins or registration
func processUser(c *gin.Context) {
	var userForm UserForm
	c.Bind(&userForm)

	var sessionID string
	var err error
	if userForm.Action == loginAction {
		sessionID, err = login(userForm.Username, userForm.Password)
	} else if userForm.Action == registerAction {
		sessionID, err = register(userForm.Username, userForm.Password)
	} else {
		err = errors.New("Invalid user action.")
	}

	if err != nil {
		invalidHandler(c, http.StatusBadRequest, err)
		return
	}

	//Give the user a session
	session := sessions.Default(c)
	clearSession(&session)
	session.Set(sessionSessionID, sessionID)
	session.Set(sessionDisplayName, userForm.Username)
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

//logout the current user
func logout(c *gin.Context) {
	//Give the user a session
	session := sessions.Default(c)
	clearSession(&session)

	c.Redirect(http.StatusFound, "/")
}
