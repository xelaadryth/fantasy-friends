package controller

import (
	"errors"
	"log"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/xelaadryth/fantasy-friends/database"
)

func validateSession(session *sessions.Session) error {
	sessionID := (*session).Get(sessionSessionID)
	sessionIDString, ok := sessionID.(string)
	if !ok || sessionIDString == "" {
		return errors.New("Malformed or missing session ID.")
	}
	_, err := database.GetUserIDFromSession(sessionIDString)
	if err != nil {
		return errors.New("Invalid or expired session.")
	}

	return nil
}

//clearSession and save it
func clearSession(session *sessions.Session) error {
	sessionID := (*session).Get(sessionSessionID)
	if sessionIDString, ok := sessionID.(string); ok && sessionIDString != "" {
		err := database.DeleteSession(sessionIDString)
		if err != nil {
			log.Println(err)
		}
	}
	(*session).Clear()
	(*session).Save()
	return nil
}
