package uwho

import (
	"net/http"
	"net/url"
)

// TODO: Probably need a function to collect errors and inform user/other people
type Hook func(*userinfo, http.ResponseWriter, *http.Request) error

func (c *coordinator) AttachServerSessionManager() {
}

func (c *coordinator) AttachClientSessionManager() {
}

func (c *coordinator) AddIdentifier(ident Identifier) {
	c.identifiers = append(c.identifiers, ident)
	last := len(c.identifiers) - 1
	c.SetHooks(&c.hooks.loggedIn, c.identifiers[last].GetLoggedOutHooks())
	c.SetHooks(&c.hooks.loggedOut, c.identifiers[last].GetLoggedInHooks())
	c.SetHooks(&c.hooks.authorized, c.identifiers[last].GetAuthorizedHooks())
	c.SetHooks(&c.hooks.aboutToLoad, c.identifiers[last].GetaboutToLoadHooks())
}

func (c *coordinator) SetHooks(existingHooks *[]*Hook, newHooks []*Hook) {
	*existingHooks = append(*existingHooks, newHooks...)
}

func (c *coordinator) CallHooks(hooks []*Hook, userinfo *userinfo, w http.ResponseWriter, r *http.Request) {
	for _, hook := range hooks {
		if err := (*hook)(userinfo, w, r); err != nil {
			defaultLogger.Error(err.Error())
		}
	}
}

func New(desiredResource, loginResult, accessDenied, logoutResult, expiredResult http.Handler,
	loginEndpoint, logoutEndpoint string) coordinator {
	loginEndpointParsed, err := url.Parse(loginEndpoint)
	logoutEndpointParsed, err := url.Parse(logoutEndpoint)
	if err != nil {
		panic(err.Error())
	}
	return coordinator{
		desiredResource: desiredResource,
		loginResult:     loginResult,
		accessDenied:    accessDenied,
		logoutResult:    logoutResult,
		expiredResult:   expiredResult,
		loginEndpoint:   loginEndpointParsed,
		logoutEndpoint:  logoutEndpointParsed}
}

type coordinator struct { // I think everything can be lower case, force initialization
	desiredResource http.Handler
	loginResult     http.Handler
	accessDenied    http.Handler
	logoutResult    http.Handler
	expiredResult   http.Handler
	loginEndpoint   *url.URL
	logoutEndpoint  *url.URL
	hooks           struct {
		loggedOut   [](*Hook)
		loggedIn    [](*Hook)
		authorized  [](*Hook)
		aboutToLoad [](*Hook)
	}
	CheckAuthorization *Hook
	identifiers        []Identifier
	// clientSessionManager
	// serverSessionManager
}

func (c *coordinator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userinfo := newUserinfo()

	if c.checkLogout(w, r) {
		userinfo.ReconcileStatus(LOGGEDOUT)
		c.CallHooks(c.hooks.loggedOut, userinfo, w, r)
		c.logoutResult.ServeHTTP(w, r)
		return
	}

	// CHECK USERSESSIONS // MAX on userinfo
	// CHECK SERVERSESSIONS //
	if userinfo.IsStatus(KNOWN) {
		(*c.CheckAuthorization)(userinfo, w, r)
	}
	if userinfo.IsStatus(AUTHORIZED) {
		c.CallHooks(c.hooks.authorized, userinfo, w, r)
		c.CallHooks(c.hooks.aboutToLoad, userinfo, w, r)
		c.desiredResource.ServeHTTP(w, r)
		return
	}

	if c.checkLogin(userinfo, w, r) {
		userinfo.NewSession()
		c.CallHooks(c.hooks.loggedIn, userinfo, w, r)
		c.CallHooks(c.hooks.aboutToLoad, userinfo, w, r)
		c.loginResult.ServeHTTP(w, r)
		return
	}

	if userinfo.IsStatus(LOGGEDOUT) || userinfo.IsStatus(EXPIRED) {
		c.CallHooks(c.hooks.loggedOut, userinfo, w, r)
		c.CallHooks(c.hooks.aboutToLoad, userinfo, w, r)
		if userinfo.IsStatus(EXPIRED) {
			c.expiredResult.ServeHTTP(w, r)
		} else {
			c.logoutResult.ServeHTTP(w, r)
		}
		return
	}

	c.CallHooks(c.hooks.aboutToLoad, userinfo, w, r)
	c.accessDenied.ServeHTTP(w, r)
}

func (c *coordinator) checkLogin(userinfo *userinfo, w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path == c.loginEndpoint.Path { // I want to do URL comparisons TODO
		for _, identifier := range c.identifiers {
			opinion, userinfoData := identifier.VerifyCredentials(w, r)
			if opinion == KNOWN {
				userinfo.ReconcileStatus(KNOWN)
				userinfo.Append(userinfoData)
			} else if opinion == SPOKEN {
				// TODO: Identifier trying to hijack whole process
			} else if opinion != UNKNOWN {
				defaultLogger.Error("An identifier is returning a strange user status: " + string(opinion))
			}
		}
	}
	return userinfo.IsStatus(KNOWN)
}

func (c *coordinator) checkLogout(w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path == c.logoutEndpoint.Path { // I want to do URL comparisons TODO
		return true
	}
	return false
}
