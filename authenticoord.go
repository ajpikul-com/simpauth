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
	hooks           struct { // Kinda broken because they don't take user states
		loggedOut   []Hook
		loggedIn    []Hook
		authorized  []Hook
		aboutToLoad []Hook
	}
	stateFactory Factory
}

func (c *coordinator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defaultLogger.Info("Serving HTTP from " + r.URL.Path)
	userStatus := NewUserStatus()
	userState := c.stateFactory.New()

	// Try to read the session
	opinion := c.sessionManager.ReadSession(userState, w, r)
	userStatus.ReconcileStatus(opinion)

	// Check to see if user is logging out
	if c.checkLogout(w, r) {
		defaultLogger.Info(r.URL.Path + ": We're about to logout")
		c.sessionManager.EndSession(userState, w, r)
		userState.DeleteState()
		userStatus.ReconcileStatus(LOGGEDOUT)
		c.CallHooks(c.hooks.loggedOut, w, r)
		c.logoutResult.ServeHTTP(w, r) // TODO, if we stay on /logout, it's a problem, logout should always move, don't stay on end points
		return
	}

	// Found a session
	if userStatus.IsStatus(KNOWN) {
		defaultLogger.Info(r.URL.Path + ": KNOWN, attempting to read data and authorize user")
		userStatus.ReconcileStatus(userState.AuthorizeUser(w, r))
		defaultLogger.Info(userStatus.StatusStr())
	}

	// User is authorized
	if userStatus.IsStatus(AUTHORIZED) {
		defaultLogger.Info(r.URL.Path + ": We are freshly authorized")
		defaultLogger.Info(userStatus.StatusStr())
		c.CallHooks(c.hooks.authorized, w, r)
		c.CallHooks(c.hooks.aboutToLoad, w, r)
		// If we want to login again (ie multiple logins), should we hijack here?
		c.desiredResource.ServeHTTP(w, r)
		return
	}

	// See if we're trying to login
	if c.checkLogin(userState, w, r) {
		defaultLogger.Info(r.URL.Path + ": checkLogin returned true")
		if err := userState.InitState(); err == ErrStateExists || err == ErrSessionExists {
			defaultLogger.Info("Logging in while Logged in?")
			// TODO I don't know where this should redirect to
			// I feel like endpoints should naturally kick people away if they are not valid
		} else if err != nil {
			// Not sure what to do here, InitState failed, so login should too TODO
			defaultLogger.Error(err.Error())
		} else {
			c.sessionManager.NewSession(userState, w, r) // May have to check for ErrSessionExists
			c.sessionManager.UpdateSession(userState, w, r)
			userStatus.ReconcileStatus(KNOWN)
		}
		c.CallHooks(c.hooks.loggedIn, w, r)
		c.CallHooks(c.hooks.aboutToLoad, w, r)
		c.loginResult.ServeHTTP(w, r)
		return
	}

	// Not authorized, known or unkown, expired
	defaultLogger.Info(r.URL.Path + " but " + userStatus.StatusStr() + " so DENIED")
	c.CallHooks(c.hooks.aboutToLoad, w, r)
	c.accessDenied.ServeHTTP(w, r)
}

func (c *coordinator) checkLogin(userState ReqByCoord, w http.ResponseWriter, r *http.Request) bool {
	loggedIn := false
	if r.URL.Path == c.loginEndpoint.Path {
		defaultLogger.Info("Equal paths")
		defaultLogger.Info(r.URL.Path)
		defaultLogger.Info(c.loginEndpoint.Path)
		for _, identifier := range c.identifiers {
			opinion := identifier.VerifyCredentials(userState, w, r)
			if opinion == KNOWN {
				loggedIn = true
				defaultLogger.Info("Found a user.")
			} else if opinion == SPOKEN {
				// TODO: Identifier trying to hijack whole process
			} else if opinion != UNKNOWN {
				defaultLogger.Error("An identifier is returning a strange user status: " + string(opinion))
			}
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
