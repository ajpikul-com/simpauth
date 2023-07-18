package main

type UserState struct {
	session string
	valid   bool
	name    string
}

func New() *UserState {
	return &UserState{}
}

func (u *UserState) VerifyUser() bool { // This signature always stays the same, you must have it
	return u.session != "" && u.name != ""
}

func (u *UserState) AcceptData(s string, n string) { // This signature will change depending on module A
	u.session = s
	u.name = n
}

func (u *UserState) ReadSession(session string) bool { // This signature will change depending on module B
	// Am I reading a valid session? // (Maybe user messed something up w/ manual changes, if so, they have no session!)
	// return false (let the handler delete that data)
	// Did I already read a session? valid should be false (How would you have, this happens before login, and it only happens once!)
	u.session = session
	return true
}

// This happens on refresh and login
func (u *UserState) MarkSession() string { // This signature will change depending on module B
	return u.session + u.name // This is just whatevers stored in the session
}

func (u *UserState) InitSession() bool { // This signature will change depending on module B
	// We've just logged in!
	u.session = "newSession"
	// Is there minimum information I need in order to do this
	return u.MarkSession() // Or maybe mark will be called seperately
}

func (u *UserState) DeleteSession() {
	// Just Delete My Own Data, We Are Now Logged Out (It shouldn't matter anyway but just in case)
}

