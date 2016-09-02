package controller

import (
	"errors"
	"log"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/xelaadryth/fantasy-friends/database"
)

func validateSession(session *sessions.Session) (int64, error) {
	sessionID := (*session).Get(sessionSessionID)
	sessionIDString, ok := sessionID.(string)
	if !ok || sessionIDString == "" {
		return 0, errors.New("Malformed or missing session ID.")
	}
	userID, err := database.GetUserIDFromSession(sessionIDString)
	if err != nil {
		return 0, errors.New("Invalid or expired session.")
	}

	return userID, nil
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
