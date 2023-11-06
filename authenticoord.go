package uwho

import (
	"net/http"
	"net/url"
)

type Coordinator struct {
	identifiers     []Identifier
	sessionManager  SessionManager
	DesiredResource http.Handler
	loginEndpoint   *url.URL
	logoutEndpoint  *url.URL
	stateFactory    Factory
	optionOverride  func(w http.ResponseWriter, r *http.Request) bool
}

// Have your handler return true if you're done talking to the user
func (c *Coordinator) OverrideOPTION(override func(w http.ResponseWriter, r *http.Request) bool) {
	c.optionOverride = override
}

// The Clone function does not deep copy. It's only purpose is to allow you to use the same uwho coordinator with a different desired resource. None of the other members are accessable anyway.
func (c *Coordinator) Clone(newResource http.Handler) *Coordinator {
	clone := new(Coordinator)
	*clone = *c
	clone.DesiredResource = newResource
	// Interface function escaping me here
	return clone
}

func (c *Coordinator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defaultLogger.Debug("**")
	defaultLogger.Debug("Coordinator start")
	defaultLogger.Debug("Serving HTTP from " + r.URL.Path + " with method " + r.Method)
	if c.optionOverride != nil && r.Method == "OPTIONS" {
		defaultLogger.Debug("We're in option!")
		if c.optionOverride(w, r) {
			defaultLogger.Debug("Override requested we not continue")
			return
		}
	}
	userState := c.stateFactory.New()

	// Read Session
	c.sessionManager.ReadSession(userState, w, r)

	// Checking Endpoints
	stateChange := false
	if r.URL.Path == c.logoutEndpoint.Path {
		defaultLogger.Debug(r.URL.Path + ": We're about to logout")
		userState.LogOut(w, r)
		stateChange = true
	} else if r.URL.Path == c.loginEndpoint.Path {
		defaultLogger.Debug("Trying to login")
		stateChange = true
		for _, identifier := range c.identifiers {
			if !userState.IsLoginAllowed(w, r) {
				break
			}
			identifier.VerifyCredentials(userState, w, r)
		}
	}

	// Finished Learning about User and Write Session
	userState.OtherStateAction(w, r)
	c.sessionManager.UpdateSession(userState, w, r)

	if stateChange {
		userState.ChangeState(w, r)
		return
	}

	if userState.IsUserAuthorized(w, r) {
		defaultLogger.Debug(r.URL.Path + ": We are freshly authorized")
		c.DesiredResource.ServeHTTP(w, r)
	}
}
