package usersessioncookie

import (
	"time"

	"github.com/ajpikul-com/uwho"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

type CookieSessionManager struct {
	id     string
	expiry time.Duration
	signer ssh.Signer
	domain string
	path   string
}

func (c *CookieSessionManager) SetID(id string) {
	c.id = id
}

func (c *CookieSessionManager) TestInterface(stateManager uwho.ReqByCoord) {
	if _, ok := stateManager.(ReqBySess); !ok {
		panic("State manager doesn't satisfied required interface")
	}
}

func New(domain string, path string, signer ssh.Signer) *CookieSessionManager {
	return &CookieSessionManager{
		id:     uuid.New().String(),
		signer: signer,
		domain: domain,
		path:   path,
	}
}
