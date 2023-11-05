package uwho

import (
	"net/url"
)

func (c *Coordinator) AttachSessionManager(m SessionManager) {
	m.TestInterface(c.stateFactory.New())
	c.sessionManager = m
}

func (c *Coordinator) AddIdentifier(ident Identifier) {
	ident.TestInterface(c.stateFactory.New())
	c.identifiers = append(c.identifiers, ident)
}

func New(loginEndpoint, logoutEndpoint string, factory Factory) *Coordinator {
	loginEndpointParsed, err := url.Parse(loginEndpoint)
	logoutEndpointParsed, err := url.Parse(logoutEndpoint)
	if err != nil {
		panic(err.Error())
	}

	return &Coordinator{
		loginEndpoint:  loginEndpointParsed,
		logoutEndpoint: logoutEndpointParsed,
		stateFactory:   factory,
	}
}
