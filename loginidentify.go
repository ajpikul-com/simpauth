package uwho

import (
	"net/http"
)

type Identifier interface {
	VerifyCredentials(http.ResponseWriter, *http.Request) (UserStatus, interface{})
	GetLoggedOutHooks() []*Hook
	GetLoggedInHooks() []*Hook
	GetAuthorizedHooks() []*Hook
	GetaboutToLoadHooks() []*Hook
}
