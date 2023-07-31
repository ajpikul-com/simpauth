package usersessioncookie

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

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
	id     uuid.UUID
	expiry time.Duration
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
func New(expiry time.Duration) *CookieSessionManager {
	return &CookieSessionManager{
		id:     uuid.New(),
		expiry: expiry,
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
		t, err := time.Parse(time.RFC3339, splitValue[1])
		if err != nil {
			c.EndSession(userStateCoord, w, r)
			return false
		}
		if c.expiry != 0 && time.Now().After(t.Add(c.expiry)) {
			c.EndSession(userStateCoord, w, r)
			expired = true
		}
		dataBits, err := base64.StdEncoding.DecodeString(splitValue[0])
		if err != nil {
			defaultLogger.Error(err.Error())
			c.EndSession(userStateCoord, w, r)
			return false
		}
		data := string(dataBits[:])
		defaultLogger.Info("Readsession captured string: " + data)
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
		http.SetCookie(w, &http.Cookie{
			Name:  c.id.String(),
			Value: base64.StdEncoding.EncodeToString(append([]byte(userState.StateToSession()+"&"), t...)),
			Path:  "/", // Maybe we should be setting this when we initialize it? Not sure how it really effects behavior
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
