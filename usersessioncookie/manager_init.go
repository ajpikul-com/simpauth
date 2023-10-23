package usersessioncookie

import (
	"os"
	"time"

	"github.com/ajpikul-com/uwho"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

type CookieSessionManager struct {
	id      string
	expiry  time.Duration
	private ssh.Signer
	domain  string
	path    string
}

func (c *CookieSessionManager) SetID(id string) {
	c.id = id
}

func (c *CookieSessionManager) TestInterface(stateManager uwho.ReqByCoord) {
	if _, ok := stateManager.(ReqBySess); !ok {
		panic("State manager doesn't satisfied required interface")
	}
}

func New(domain string, path string, key string) *CookieSessionManager {
	privateBytes, err := os.ReadFile(key)
	if err != nil {
		panic(err.Error())
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic(err.Error())
	}

	return &CookieSessionManager{
		id:      uuid.New().String(),
		private: private,
		domain:  domain,
		path:    path,
	}
}
