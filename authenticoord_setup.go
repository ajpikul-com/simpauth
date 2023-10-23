package uwho

import (
	"net/url"
)

func (c *coordinator) AttachSessionManager(m SessionManager) {
	m.TestInterface(c.stateFactory.New())
	c.sessionManager = m
}

func (c *coordinator) AddIdentifier(ident Identifier) {
	ident.TestInterface(c.stateFactory.New())
	c.identifiers = append(c.identifiers, ident)
}

func New(loginEndpoint, logoutEndpoint string, factory Factory) coordinator {
	loginEndpointParsed, err := url.Parse(loginEndpoint)
	logoutEndpointParsed, err := url.Parse(logoutEndpoint)
	if err != nil {
		panic(err.Error())
	}

	return coordinator{
		loginEndpoint:  loginEndpointParsed,
		logoutEndpoint: logoutEndpointParsed,
		stateFactory:   factory,
	}
}
