package uwho

import (
	"errors"
	"net/http"
)

type Hook *func(http.ResponseWriter, *http.Request) error

type Module interface {
	GetLoggedOutHooks() []Hook
	GetLoggedInHooks() []Hook
	GetAuthorizedHooks() []Hook
	GetAboutToLoadHooks() []Hook
}
type Identifier interface {
	Module
	VerifyCredentials(http.ResponseWriter, *http.Request) UserStatus
}

var ErrSessionExists error = errors.New("Session already exist")

type SessionManager interface {
	Module
	ReadSession(http.ResponseWriter, *http.Request) (string, UserStatus)
	MarkSession(string, http.ResponseWriter, *http.Request)
	EndSession(http.ResponseWriter, *http.Request)
}

type AppUserinfo interface {
	LogOut()
	SessionString() string
	SessionDestring(string)
	AuthorizeUser(w http.ResponseWriter, r *http.Request) UserStatus
	InitSession() error
}
