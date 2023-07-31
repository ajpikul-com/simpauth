package uwho

import (
	"net/http"
	"net/url"
)

// TODO: Probably need a function to collect errors and inform user/other people

type coordinator struct {
	identifiers     []Identifier
	sessionManager  SessionManager
	desiredResource http.Handler
	loginResult     http.Handler
	accessDenied    http.Handler
	logoutResult    http.Handler
	loginEndpoint   *url.URL
	logoutEndpoint  *url.URL
	Hooks           struct { // Kinda broken because they don't take user states
		LoggedOut   []Hook
		LoggedIn    []Hook
		Authorized  []Hook
		AboutToLoad []Hook
	}
	stateFactory Factory
}

func (c *coordinator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defaultLogger.Info("Serving HTTP from " + r.URL.Path)
	userStatus := NewUserStatus()
	userState := c.stateFactory.New()

	// Try to read the session
	if c.sessionManager.ReadSession(userState, w, r) {
		userStatus.ReconcileStatus(KNOWN)
	}

	// KNOWN OR UNKNOWN

	// Checking Endpoint
	if c.checkLogout(w, r) {
		defaultLogger.Info(r.URL.Path + ": We're about to logout")
		c.sessionManager.EndSession(userState, w, r)
		userState.DeleteState()
		userStatus.ReconcileStatus(LOGGEDOUT)
		c.CallHooks(c.Hooks.LoggedOut, userState, w, r)
		c.logoutResult.ServeHTTP(w, r)
		return
	} else if c.checkLogin(userState, w, r) {
		defaultLogger.Info(r.URL.Path + ": checkLogin returned true")
		if userStatus.IsStatus(UNKNOWN) {
			err := userState.InitState()
			if err != nil {
				defaultLogger.Error(err.Error()) // Login Likely Failed
			} else {
				c.sessionManager.NewSession(userState, w, r)    // May have to check for ErrSessionExists
				c.sessionManager.UpdateSession(userState, w, r) // May have to check for Errors
				userStatus.ReconcileStatus(KNOWN)
				c.CallHooks(c.Hooks.LoggedIn, userState, w, r)
			}
		}
		c.CallHooks(c.Hooks.AboutToLoad, userState, w, r)
		c.loginResult.ServeHTTP(w, r)
		return
	} else if r.URL.Path == c.loginEndpoint.Path {
		c.CallHooks(c.Hooks.AboutToLoad, userState, w, r)
		c.accessDenied.ServeHTTP(w, r)
		return
	}

	// KNOWN OR UNKNOWN (LOGGEDOUT RETURNED)

	if userStatus.IsStatus(KNOWN) {
		defaultLogger.Info(r.URL.Path + ": KNOWN, attempting to read data and authorize user")
		if userState.AuthorizeUser(w, r) {
			userStatus.ReconcileStatus(AUTHORIZED)
		}
		defaultLogger.Info(userStatus.StatusStr())
	}

	// KNOWN, UNKNOWN, OR AUTHORIZED (MAYBE UNAUTHORIZED)

	// User is authorized
	if userStatus.IsStatus(AUTHORIZED) {
		defaultLogger.Info(r.URL.Path + ": We are freshly authorized")
		defaultLogger.Info(userStatus.StatusStr())
		c.CallHooks(c.Hooks.Authorized, userState, w, r)
		c.CallHooks(c.Hooks.AboutToLoad, userState, w, r)
		c.desiredResource.ServeHTTP(w, r)
		return
	}

	c.CallHooks(c.Hooks.AboutToLoad, userState, w, r)
	c.accessDenied.ServeHTTP(w, r)
}

func (c *coordinator) checkLogin(userState ReqByCoord, w http.ResponseWriter, r *http.Request) bool {
	loggedIn := false
	if r.URL.Path == c.loginEndpoint.Path {
		defaultLogger.Info("Equal paths")
		defaultLogger.Info(r.URL.Path)
		defaultLogger.Info(c.loginEndpoint.Path)
		for _, identifier := range c.identifiers {
			loggedIn = identifier.VerifyCredentials(userState, w, r)
		}
	}
	return loggedIn
}

func (c *coordinator) checkLogout(w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path == c.logoutEndpoint.Path { // I want to do URL comparisons TODO
		return true
	}
	return false
}
