package uwho

import (
	"errors"
	"net/http"
)

type Module interface {
	TestInterface(ReqByCoord)
}

type Identifier interface {
	Module
	VerifyCredentials(ReqByCoord, http.ResponseWriter, *http.Request)
}

var ErrSessionExists error = errors.New("Session already exist")

type SessionManager interface {
	Module
	ReadSession(ReqByCoord, http.ResponseWriter, *http.Request)
	UpdateSession(ReqByCoord, http.ResponseWriter, *http.Request)
}

// INTERFACES THAT MUST BE SATISFIED BY USER-DEVELOPER IN THEIR STATE MANAGER
type Factory interface {
	New() ReqByCoord
}

type ReqByCoord interface {
	LogOut(w http.ResponseWriter, r *http.Request)
	IsLoginAllowed(w http.ResponseWriter, r *http.Request) bool
	OtherStateAction(w http.ResponseWriter, r *http.Request)
	ChangeState(w http.ResponseWriter, r *http.Request)
	IsUserAuthorized(w http.ResponseWriter, r *http.Request) bool
}
