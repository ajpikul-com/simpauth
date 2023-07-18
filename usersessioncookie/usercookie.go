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

type CookieSessionManager struct {
	id uuid.UUID
}

func (g *CookieSessionManager) GetLoggedOutHooks() []uwho.Hook   { return nil }
func (g *CookieSessionManager) GetLoggedInHooks() []uwho.Hook    { return nil }
func (g *CookieSessionManager) GetAuthorizedHooks() []uwho.Hook  { return nil }
func (g *CookieSessionManager) GetAboutToLoadHooks() []uwho.Hook { return nil }

func New() *CookieSessionManager {
	return &CookieSessionManager{
		id: uuid.New(),
	}
}

// TODO: needs some encryption or signing here

func (m *CookieSessionManager) ReadSession(w http.ResponseWriter, r *http.Request) (string, uwho.UserStatus) {
	cookie, err := r.Cookie(m.id.String())
	if err == http.ErrNoCookie {
		return "", uwho.UNKNOWN
	} else if err == nil {
		dataBits, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err != nil {
			defaultLogger.Error(err.Error())
			return "", uwho.UNKNOWN
		}
		data := string(dataBits[:])
		defaultLogger.Info("Readsession captured string: " + data)
		return data, uwho.KNOWN
	} else {
		defaultLogger.Error(err.Error())
		return "", uwho.UNKNOWN
	}
}

func (m *CookieSessionManager) MarkSession(data string, w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  m.id.String(),
		Value: base64.StdEncoding.EncodeToString([]byte(data)),
		Path:  "/", // Maybe we should be setting this when we initialize it? Not sure how it really effects behavior
	})
}

func (m *CookieSessionManager) EndSession(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   m.id.String(),
		Value:  "",
		MaxAge: -1,
		Path:   "/", // Maybe we should be setting this when we initialize it? Not sure how it really effects behavior
	})
}

// TODO NEED A WAY TO LOGOUT, OR DELETE COOKIES, WHICH WE CALL IF WE'RE EXPIRED!
