package uwho

import (
	"net/http"
	"net/url"
)

func (c *coordinator) AttachSessionManager(m SessionManager) {
	m.TestInterface(c.stateFactory.New())
	c.sessionManager = m
	c.SetHooks(&c.Hooks.LoggedIn, c.sessionManager.GetLoggedOutHooks())
	c.SetHooks(&c.Hooks.LoggedOut, c.sessionManager.GetLoggedInHooks())
	c.SetHooks(&c.Hooks.Authorized, c.sessionManager.GetAuthorizedHooks())
	c.SetHooks(&c.Hooks.AboutToLoad, c.sessionManager.GetAboutToLoadHooks())
}

func (c *coordinator) AddIdentifier(ident Identifier) {
	ident.TestInterface(c.stateFactory.New())
	c.identifiers = append(c.identifiers, ident)
	last := len(c.identifiers) - 1
	c.SetHooks(&c.Hooks.LoggedIn, c.identifiers[last].GetLoggedOutHooks())
	c.SetHooks(&c.Hooks.LoggedOut, c.identifiers[last].GetLoggedInHooks())
	c.SetHooks(&c.Hooks.Authorized, c.identifiers[last].GetAuthorizedHooks())
	c.SetHooks(&c.Hooks.AboutToLoad, c.identifiers[last].GetAboutToLoadHooks())
}

func (c *coordinator) SetHooks(existingHooks *[]Hook, newHooks []Hook) {
	*existingHooks = append(*existingHooks, newHooks...)
}

// TODO Unsure how it calls user data
func (c *coordinator) CallHooks(hooks []Hook, state ReqByCoord, w http.ResponseWriter, r *http.Request) {
	for _, hook := range hooks {
		if err := (*hook)(state, w, r); err != nil {
			defaultLogger.Error(err.Error())
		}
	}
}

// TODO, this is going to change with refactor, since now we take a factory
func New(desiredResource, loginResult, accessDenied, logoutResult http.Handler,
	loginEndpoint, logoutEndpoint string,
	factory Factory) coordinator {

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
		loginEndpoint:   loginEndpointParsed,
		logoutEndpoint:  logoutEndpointParsed,
		stateFactory:    factory,
	}
}
