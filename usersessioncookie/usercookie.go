package usersessioncookie

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/ajpikul-com/ilog"
	"github.com/ajpikul-com/uwho"
	"github.com/google/uuid"
)

var defaultLogger ilog.LoggerInterface

func init() {
	if defaultLogger == nil {
		defaultLogger = new(ilog.EmptyLogger)
	}
}

func SetDefaultLogger(newLogger ilog.LoggerInterface) {
	defaultLogger = newLogger
	defaultLogger.Info("Default Logger Set")
}

type ReqBySess interface {
	StateToSession() string
	SessionToState(string, bool) bool // Please return false if you don't think the session is legit
}

type CookieSessionManager struct {
	id      uuid.UUID
	expiry  time.Duration
	private ssh.Signer
	domain  string
	path    string
}

func (c *CookieSessionManager) GetLoggedOutHooks() []uwho.Hook   { return nil }
func (c *CookieSessionManager) GetLoggedInHooks() []uwho.Hook    { return nil }
func (c *CookieSessionManager) GetAuthorizedHooks() []uwho.Hook  { return nil }
func (c *CookieSessionManager) GetAboutToLoadHooks() []uwho.Hook { return nil }
func (c *CookieSessionManager) TestInterface(stateManager uwho.ReqByCoord) {
	if _, ok := stateManager.(ReqBySess); !ok {
		panic("State manager doesn't satisfied required interface")
	}
}
func New(domain string, path string, expiry time.Duration, key string) *CookieSessionManager {
	privateBytes, err := os.ReadFile(key)
	if err != nil {
		panic(err.Error())
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic(err.Error())
	}

	return &CookieSessionManager{
		id:      uuid.New(),
		expiry:  expiry,
		private: private,
		domain:  domain,
		path:    path,
	}
}

func (c *CookieSessionManager) NewSession(userStateCoord uwho.ReqByCoord, w http.ResponseWriter, r *http.Request) {
	// No prep todo with new session, we can just update session, at least on this implementation
}

func (c *CookieSessionManager) ReadSession(userStateCoord uwho.ReqByCoord, w http.ResponseWriter, r *http.Request) bool {
	expired := false
	cookie, err := r.Cookie(c.id.String())
	if err == http.ErrNoCookie {
		c.EndSession(userStateCoord, w, r)
		return false
	} else if err == nil {
		splitValue := strings.Split(cookie.Value, "&")
		blob, err := base64.StdEncoding.DecodeString(splitValue[3])
		if err != nil {
			c.EndSession(userStateCoord, w, r)
			return false
		}
		signature := &ssh.Signature{
			Format: splitValue[2],
			Blob:   blob,
		}
		if len(splitValue) == 5 {
			rest, err := base64.StdEncoding.DecodeString(splitValue[4])
			if err != nil {
				c.EndSession(userStateCoord, w, r)
				return false
			}
			signature.Rest = rest
			defaultLogger.Error("Read a sig string:")
			defaultLogger.Error(splitValue[2] + "&" + splitValue[3] + "&" + splitValue[4])
		} else {
			defaultLogger.Error("Read a sig string:")
			defaultLogger.Error(splitValue[2] + "&" + splitValue[3])
		}
		dataBits, err := base64.StdEncoding.DecodeString(splitValue[0])
		if err != nil {
			defaultLogger.Error(err.Error())
			c.EndSession(userStateCoord, w, r)
			return false
		}
		err = c.private.PublicKey().Verify(dataBits, signature)
		if err != nil {
			defaultLogger.Debug(err.Error())
			c.EndSession(userStateCoord, w, r)
			return false
		}
		t, err := time.Parse(time.RFC3339, splitValue[1])
		if err != nil {
			c.EndSession(userStateCoord, w, r)
			return false
		}
		if c.expiry != 0 && time.Now().After(t.Add(c.expiry)) {
			c.EndSession(userStateCoord, w, r)
			expired = true
		}
		data := string(dataBits[:])
		defaultLogger.Debug("Readsession captured string: " + data)
		if userState, ok := userStateCoord.(ReqBySess); ok {
			ok = userState.SessionToState(data, expired)
			if !ok || expired {
				c.EndSession(userStateCoord, w, r)
				return false
			}
		} else {
			c.EndSession(userStateCoord, w, r)
			return false
		}
		c.UpdateSession(userStateCoord, w, r)
		return true
	} else {
		defaultLogger.Error(err.Error())
		c.EndSession(userStateCoord, w, r)
		return false
	}
}

func (c *CookieSessionManager) UpdateSession(userStateCoord uwho.ReqByCoord, w http.ResponseWriter, r *http.Request) {
	if userState, ok := userStateCoord.(ReqBySess); ok {
		t, _ := time.Now().MarshalText()
		bytes := []byte(userState.StateToSession())
		signature, _ := c.private.Sign(rand.Reader, bytes)
		sigString := signature.Format + "&" + base64.StdEncoding.EncodeToString(signature.Blob)
		if len(signature.Rest) != 0 {
			sigString += "&" + base64.StdEncoding.EncodeToString(signature.Rest)
		}
		defaultLogger.Error("Generated up a sig string:")
		defaultLogger.Error(sigString)
		value := base64.StdEncoding.EncodeToString(bytes) + "&" + string(t[:]) + "&" + sigString
		http.SetCookie(w, &http.Cookie{
			Name:   c.id.String(),
			Value:  value,
			Domain: c.domain,
			Path:   c.path,
		})
	}
}

func (c *CookieSessionManager) EndSession(userStateCoord uwho.ReqByCoord, w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   c.id.String(),
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})
}
