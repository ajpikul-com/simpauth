package uwho

import (
	"errors"
	"net/http"
)

// INTERFACES THAT MUST BE SATISFIED BY MODULES
type Hook *func(ReqByCoord, http.ResponseWriter, *http.Request) error

type Module interface {
	TestInterface(ReqByCoord)
	GetLoggedOutHooks() []Hook
	GetLoggedInHooks() []Hook
	GetAuthorizedHooks() []Hook
	GetAboutToLoadHooks() []Hook
}

type Identifier interface {
	Module
	VerifyCredentials(ReqByCoord, http.ResponseWriter, *http.Request) bool
}

var ErrSessionExists error = errors.New("Session already exist")
var ErrStateExists error = errors.New("State already exist")
var ErrNoCredential error = errors.New("Login Failed")

type SessionManager interface {
	Module
	NewSession(ReqByCoord, http.ResponseWriter, *http.Request)
	ReadSession(ReqByCoord, http.ResponseWriter, *http.Request) bool
	UpdateSession(ReqByCoord, http.ResponseWriter, *http.Request)
	EndSession(ReqByCoord, http.ResponseWriter, *http.Request)
}

// INTERFACES THAT MUST BE SATISFIED BY USER-DEVELOPER IN THEIR STATE MANAGER
type Factory interface {
	New() ReqByCoord
}

type ReqByCoord interface {
	AuthorizeUser(w http.ResponseWriter, r *http.Request) bool
	InitState() error
	DeleteState()
}
