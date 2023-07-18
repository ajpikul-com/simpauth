package main

type ReqByCoord interface {
	AuthorizeUser() bool
	InitState()
	DeleteState()
	AddStateListener(func()) // Your structure can embed a structure that does this	from uwho
	StateChangedHook()       // Your structure can embed a structure that does this from uwho
}
type ReqByIdent interface {
	AcceptData(s string, n string)
}
type ReqBySess interface {
	StateToSession() string
	SessionToState(session string, name string) bool
}

type Coordinator struct{}
type Identifier struct{}
type SessionManager struct{}
// Problem: we will be creating new userstates, but before we start (at run time, or compile time), we want to make sure our userstate's will satisfy several interfaces.
// Our coordinator cannot accept a type. It can accept a test, though.

type UserState struct {
	session string
	valid   bool
	name    string
}

func New() *UserState {
	return &UserState{}
}

func (u *UserState) AuthorizeUser() bool { // REQUIRED BY COORDINATOR
	return u.session != "" && u.name != ""
}

func (u *UserState) AcceptData(s string, n string) { // REQUIRED BY IDENTIFIER
	u.session = s
	u.name = n
}

// CALLED BY COORDINATOR AFTER LOGIN, BEFORE INIT BY SESSION MANAGER
func (u *UserState) InitState() { // REQUIRED/DECLARED BY COORIDINATOR
	// We've just logged in! We Need to create state for the first time
	u.session = "newSession"
	u.valid = true
}

// CALLED BY COORDINATOR BUT AFTER DELETE BY SESSION MANAGER
func (u *UserState) DeleteState() { // REQUIRED/DECLARED BY COORIDINATOR
	// Just Delete My Own Data, We Are Now Logged Out (It shouldn't matter anyway but just in case)
	return u.name
}

// SESSION MANAGER HAS FOUND A SESSION
func (u *UserState) SessionToState(session string, name string) bool { // REQUIRED/DECLARED BY SESSION MANAGER
	// Am I reading a valid session? Me and the session manager have to agree. Session manager must not ignore me.
	if session == "" || name == false {
		return false // I don't have what I need, we need to make the user login again
	}
	u.session = session
	u.name = name
	return true
}

// SESSION MANAGER WANTS TO UPDATE DETAILS
func (u *UserState) StateToSession() string { // REQUIRED/DECLARED BY SESSION MANAGER
	return u.session, u.name
}

// IF USER CHANGES STATE DETAILS
func (u *UserState) AddStateListener(func()) { // REQUIRED BY COORDINATOR
	// This gets called if you or the user changes any information that should be reflected in session stoarge
}

// IF USER CHANGES STATE DETAILS
func (u *UserState) StateChangedHook() { // REQUIRED BY COORDINATOR
	// This gets called if you or the user changes any information that should be reflected in session stoarge
}

