package main

import "fmt"

type Factory interface {
	New() ReqByCoord
}
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
	StateToSession() (string, string)
	SessionToState(session string, name string) bool
}

type Coordinator struct{ factory Factory }

func (c *Coordinator) GenerateState() ReqByCoord {
	return c.factory.New()
}

type Identifier struct{}
type SessionManager struct{}

type UserStateFactory struct{}

func (uf UserStateFactory) New() ReqByCoord {
	return &UserState{}
}

type UserState struct {
	session string
	valid   bool
	name    string
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
}

// SESSION MANAGER HAS FOUND A SESSION
func (u *UserState) SessionToState(session string, name string) bool { // REQUIRED/DECLARED BY SESSION MANAGER
	// Am I reading a valid session? Me and the session manager have to agree. Session manager must not ignore me.
	if session == "" || name == "" {
		return false // I don't have what I need, we need to make the user login again
	}
	u.session = session
	u.name = name
	return true
}

// SESSION MANAGER WANTS TO UPDATE DETAILS
func (u *UserState) StateToSession() (string, string) { // REQUIRED/DECLARED BY SESSION MANAGER
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

func main() {
	stateFactory := UserStateFactory{}
	coordinator := &Coordinator{factory: stateFactory}
	actualState := coordinator.GenerateState()

	// This could be a runtime check, on confinguration
	coord, ok := actualState.(ReqByCoord)
	fmt.Println(coord, ok)
	ident, ok := actualState.(ReqByIdent)
	fmt.Println(ident, ok)
	sess, ok := actualState.(ReqBySess)
	fmt.Println(sess, ok)
	state, ok := actualState.(*UserState)
	fmt.Println(*state, ok)
	
	// So each backends guy will have to call one of these. Then each request will have to set up all its type conversions. Lets see if they're all modifying the same thing.
}

