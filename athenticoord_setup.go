package uwho

import (
	"net/http"
	"net/url"
)

func (c *coordinator) AttachSessionManager(m SessionManager) {
	c.sessionManager = m
	c.SetHooks(&c.hooks.loggedIn, c.sessionManager.GetLoggedOutHooks())
	c.SetHooks(&c.hooks.loggedOut, c.sessionManager.GetLoggedInHooks())
	c.SetHooks(&c.hooks.authorized, c.sessionManager.GetAuthorizedHooks())
	c.SetHooks(&c.hooks.aboutToLoad, c.sessionManager.GetAboutToLoadHooks())
}

func (c *coordinator) AddIdentifier(ident Identifier) {
	c.identifiers = append(c.identifiers, ident)
	last := len(c.identifiers) - 1
	c.SetHooks(&c.hooks.loggedIn, c.identifiers[last].GetLoggedOutHooks())
	c.SetHooks(&c.hooks.loggedOut, c.identifiers[last].GetLoggedInHooks())
	c.SetHooks(&c.hooks.authorized, c.identifiers[last].GetAuthorizedHooks())
	c.SetHooks(&c.hooks.aboutToLoad, c.identifiers[last].GetAboutToLoadHooks())
}

func (c *coordinator) SetHooks(existingHooks *[]Hook, newHooks []Hook) {
	*existingHooks = append(*existingHooks, newHooks...)
}

// TODO unsure how it receives user data
func (c *coordinator) CallHooks(hooks []Hook, w http.ResponseWriter, r *http.Request) {
	for _, hook := range hooks {
		if err := (*hook)(w, r); err != nil {
			defaultLogger.Error(err.Error())
		}
	}
}

func New(desiredResource, loginResult, accessDenied, logoutResult, expiredResult http.Handler,
	loginEndpoint, logoutEndpoint string,
	applicationUserinfo AppUserinfo) coordinator {

	loginEndpointParsed, err := url.Parse(loginEndpoint)
	logoutEndpointParsed, err := url.Parse(logoutEndpoint)
	if err != nil {
		panic(err.Error())
	}

	return coordinator{
		desiredResource:     desiredResource,
		loginResult:         loginResult,
		accessDenied:        accessDenied,
		logoutResult:        logoutResult,
		expiredResult:       expiredResult,
		loginEndpoint:       loginEndpointParsed,
		logoutEndpoint:      logoutEndpointParsed,
		applicationUserinfo: applicationUserinfo,
	}
}
