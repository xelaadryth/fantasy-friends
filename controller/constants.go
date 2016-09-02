package controller

const (
	sessionSession     = "session"
	sessionSessionID   = "sessionID"
	sessionDisplayName = "displayName"
	sessionNavActive   = "navActive"
	sessionTeam        = "team"
	sessionName        = "name"
	sessionTop         = "top"
	sessionJungle      = "jungle"
	sessionMid         = "mid"
	sessionBottom      = "bottom"
	sessionSupport     = "support"
)

//sessionFields that we want to hold at the top level of the session object
var sessionFields = [...]string{sessionDisplayName, sessionNavActive, sessionTeam}
