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
	// LogOut will empty out user state object, that's all. It should _not_ write a response body.
	LogOut(w http.ResponseWriter, r *http.Request)
	// IsLoginAllowed tells uwho to skip login process if, for example, we're already logged in.
	// hint: you can also use it to hint that user is trying to login. It should _not_ write a response body.
	IsLoginAllowed(w http.ResponseWriter, r *http.Request) bool
	// OtherStateAction is a hook that will be called after every other source of user information has been requested. It should _not_ write a response body.
	OtherStateAction(w http.ResponseWriter, r *http.Request)
	// ChangeState will be called if the user has tried to login (with success or not), or loggedout. It must write a response body. `uwho` provides some obvious utility functions (see README.md or utilities.go) that you can use.
	ChangeState(w http.ResponseWriter, r *http.Request)
	// IsUserAuthorized will be called after session is read, the user did not login or logout, and session is updated. It's your job to check the request and see if user is authorized. If true, user will continue to the wrapped handler, Coordinator.DesiredResource. If false, you must write a response body. Maybe redirect user to a login page?
	IsUserAuthorized(w http.ResponseWriter, r *http.Request) bool
}
