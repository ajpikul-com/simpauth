package uwho

import (
	"net/http"
	"net/url"
)

type coordinator struct {
	identifiers     []Identifier
	sessionManager  SessionManager
	desiredResource http.Handler
	loginEndpoint   *url.URL
	logoutEndpoint  *url.URL
	stateFactory    Factory
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
		c.desiredResource.ServeHTTP(w, r)
	}
}
