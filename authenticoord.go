package uwho

import (
	"net/http"
	"net/url"
)

type coordinator struct {
	identifiers     []Identifier
	sessionManager  SessionManager
	DesiredResource http.Handler
	loginEndpoint   *url.URL
	logoutEndpoint  *url.URL
	stateFactory    Factory
}

// The Clone function does not deep copy. It's only purpose is to allow you to use the same uwho coordinator with a different desired resource. None of the other members are accessable anyway.
func (c *coordinator) Clone(newResource http.Handler) *coordinator {
	clone := new(coordinator)
	*clone = *c
	clone.DesiredResource = newResource
	// Interface function escaping me here
	return clone
}

func (c *coordinator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defaultLogger.Debug("Serving HTTP from " + r.URL.Path)
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
