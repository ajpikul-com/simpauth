package uwho

import (
	"errors"
	"net/http"
)

// INTERFACES THAT MUST BE SATISFIED BY MODULES
type Hook *func(http.ResponseWriter, *http.Request) error

type Module interface {
	TestInterface(ReqByCoord)
	GetLoggedOutHooks() []Hook
	GetLoggedInHooks() []Hook
	GetAuthorizedHooks() []Hook
	GetAboutToLoadHooks() []Hook
}

type Identifier interface {
	Module
	VerifyCredentials(ReqByCoord, http.ResponseWriter, *http.Request) UserStatus
}

var ErrSessionExists error = errors.New("Session already exist")

type SessionManager interface {
	Module
	NewSession(http.ResponseWriter, *http.Request)
	ReadSession(ReqByCoord, http.ResponseWriter, *http.Request) UserStatus // The return here is effectively an error
	UpdateSession(ReqByCoord, http.ResponseWriter, *http.Request)
	EndSession(http.ResponseWriter, *http.Request)
}

// INTERFACES THAT MUST BE SATISFIED BY USER-DEVELOPER IN THEIR STATE MANAGER
type Factory interface {
	New() ReqByCoord
}

type ReqByCoord interface {
	AuthorizeUser(w http.ResponseWriter, r *http.Request) UserStatus
	InitState() error
	DeleteState()
} // Need those hooks back. SessionManager.UpdateSession needs to get registered with the hook.
