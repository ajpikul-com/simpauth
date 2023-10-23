package uwho

import (
	"net/url"
)

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
