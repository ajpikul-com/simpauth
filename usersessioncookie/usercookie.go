package usersessioncookie

import (
	"encoding/base64"
	"net/http"

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
	SessionToState(string) bool // Not sure why again
}

type CookieSessionManager struct {
	id uuid.UUID
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
func New() *CookieSessionManager {
	return &CookieSessionManager{
		id: uuid.New(),
	}
}

func (c *CookieSessionManager) NewSession(http.ResponseWriter, *http.Request) {
	// No prep todo with new session, we can just update session
}

func (c *CookieSessionManager) ReadSession(userStateCoord uwho.ReqByCoord, w http.ResponseWriter, r *http.Request) uwho.UserStatus {
	cookie, err := r.Cookie(c.id.String())
	if err == http.ErrNoCookie {
		return uwho.UNKNOWN
	} else if err == nil {
		dataBits, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err != nil {
			defaultLogger.Error(err.Error())
			return uwho.UNKNOWN
		}
		data := string(dataBits[:])
		defaultLogger.Info("Readsession captured string: " + data)
		if userState, ok := userStateCoord.(ReqBySess); ok {
			ok = userState.SessionToState(data)
			if !ok {
				return uwho.UNKNOWN
			}
		} else {
			return uwho.UNKNOWN
		}
		return uwho.KNOWN
	} else {
		defaultLogger.Error(err.Error())
		return uwho.UNKNOWN
	}
}

func (c *CookieSessionManager) UpdateSession(userStateCoord uwho.ReqByCoord, w http.ResponseWriter, r *http.Request) {
	if userState, ok := userStateCoord.(ReqBySess); ok {
		http.SetCookie(w, &http.Cookie{
			Name:  c.id.String(),
			Value: base64.StdEncoding.EncodeToString([]byte(userState.StateToSession())),
			Path:  "/", // Maybe we should be setting this when we initialize it? Not sure how it really effects behavior
		})
	}
}

func (c *CookieSessionManager) EndSession(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   c.id.String(),
		Value:  "",
		MaxAge: -1,
		Path:   "/", // Maybe we should be setting this when we initialize it? Not sure how it really effects behavior
	})
}

// TODO NEED A WAY TO LOGOUT, OR DELETE COOKIES, WHICH WE CALL IF WE'RE EXPIRED!
